package events

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"job4j_go_share_trip_contracts/internal/domain/contract/entity"
	"job4j_go_share_trip_contracts/internal/domain/contract/repository"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := postgres.Run(
		ctx,
		"postgres:16",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("password"),
	)
	if err != nil {
		log.Fatalf("start postgres container: %v", err)
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("get connection string: %v", err)
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("open sql db: %v", err)
	}

	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		pingCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		err = db.PingContext(pingCtx)
		cancel()
		if err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		log.Fatalf("database is not ready: %v", err)
	}

	if err = goose.SetDialect("postgres"); err != nil {
		log.Fatalf("set goose dialect: %v", err)
	}
	if err = goose.Up(db, "../../migrations"); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	testPool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("create pgx pool: %v", err)
	}

	code := m.Run()

	testPool.Close()
	_ = db.Close()
	_ = container.Terminate(ctx)

	os.Exit(code)
}

func newTestProcessor() *Processor {
	return NewProcessor(
		testPool,
		repository.NewUserContractRepository(),
		repository.NewProcessedEventRepository(),
	)
}

func countContracts(t *testing.T, ctx context.Context, userID uuid.UUID) int {
	t.Helper()
	var count int
	err := testPool.QueryRow(ctx, `SELECT COUNT(*) FROM contracts WHERE user_id = $1`, userID).Scan(&count)
	require.NoError(t, err)
	return count
}

func TestProcessor_CreatedCreatesDraftContract(t *testing.T) {
	ctx := context.Background()
	processor := newTestProcessor()
	driverID := uuid.New()
	companyID := uuid.New()

	err := processor.Handle(ctx, TripEvent{
		EventID:   uuid.New().String(),
		EventType: EventTypeTripCreated,
		TripID:    uuid.New().String(),
		DriverID:  driverID.String(),
		CompanyID: companyID.String(),
	})
	require.NoError(t, err)
	require.Equal(t, 1, countContracts(t, ctx, driverID))

	var status string
	var startDate sql.NullTime
	err = testPool.QueryRow(ctx,
		`SELECT status, start_date FROM contracts WHERE user_id = $1`,
		driverID,
	).Scan(&status, &startDate)
	require.NoError(t, err)
	assert.Equal(t, string(entity.StatusDraft), status)
	assert.False(t, startDate.Valid)
}

func TestProcessor_CreatedSkipsWhenOpenExists(t *testing.T) {
	ctx := context.Background()
	processor := newTestProcessor()
	driverID := uuid.New()
	companyID := uuid.New()
	existingID := uuid.New()

	_, err := testPool.Exec(ctx, `
		INSERT INTO contracts (id, company_id, user_id, status, created_at)
		VALUES ($1, $2, $3, 'draft', now())
	`, existingID, companyID, driverID)
	require.NoError(t, err)

	err = processor.Handle(ctx, TripEvent{
		EventID:   uuid.New().String(),
		EventType: EventTypeTripCreated,
		DriverID:  driverID.String(),
		CompanyID: uuid.New().String(),
	})
	require.NoError(t, err)
	require.Equal(t, 1, countContracts(t, ctx, driverID))

	var status string
	var id uuid.UUID
	err = testPool.QueryRow(ctx,
		`SELECT id, status FROM contracts WHERE user_id = $1`,
		driverID,
	).Scan(&id, &status)
	require.NoError(t, err)
	assert.Equal(t, existingID, id)
	assert.Equal(t, string(entity.StatusDraft), status)
}

func TestProcessor_CreatedCreatesNewWhenOnlyTerminated(t *testing.T) {
	ctx := context.Background()
	processor := newTestProcessor()
	driverID := uuid.New()
	companyID := uuid.New()

	_, err := testPool.Exec(ctx, `
		INSERT INTO contracts (id, company_id, user_id, status, start_date, end_date, created_at)
		VALUES ($1, $2, $3, 'terminated', now(), now(), now())
	`, uuid.New(), companyID, driverID)
	require.NoError(t, err)

	err = processor.Handle(ctx, TripEvent{
		EventID:   uuid.New().String(),
		EventType: EventTypeTripCreated,
		DriverID:  driverID.String(),
		CompanyID: companyID.String(),
	})
	require.NoError(t, err)
	require.Equal(t, 2, countContracts(t, ctx, driverID))

	var draftCount int
	err = testPool.QueryRow(ctx,
		`SELECT COUNT(*) FROM contracts WHERE user_id = $1 AND status = 'draft'`,
		driverID,
	).Scan(&draftCount)
	require.NoError(t, err)
	assert.Equal(t, 1, draftCount)
}

func TestProcessor_CreatedIsIdempotent(t *testing.T) {
	ctx := context.Background()
	processor := newTestProcessor()
	driverID := uuid.New()
	event := TripEvent{
		EventID:   uuid.New().String(),
		EventType: EventTypeTripCreated,
		DriverID:  driverID.String(),
		CompanyID: uuid.New().String(),
	}

	require.NoError(t, processor.Handle(ctx, event))
	require.NoError(t, processor.Handle(ctx, event))
	require.Equal(t, 1, countContracts(t, ctx, driverID))
}

func TestProcessor_PublishedActivatesDraftContract(t *testing.T) {
	ctx := context.Background()
	processor := newTestProcessor()
	driverID := uuid.New()
	companyID := uuid.New()
	existingID := uuid.New()
	occurredAt := time.Now().UTC().Truncate(time.Second)

	_, err := testPool.Exec(ctx, `
		INSERT INTO contracts (id, company_id, user_id, status, created_at)
		VALUES ($1, $2, $3, 'draft', now())
	`, existingID, companyID, driverID)
	require.NoError(t, err)

	err = processor.Handle(ctx, TripEvent{
		EventID:    uuid.New().String(),
		EventType:  EventTypeTripPublished,
		DriverID:   driverID.String(),
		CompanyID:  companyID.String(),
		OccurredAt: occurredAt,
	})
	require.NoError(t, err)
	require.Equal(t, 1, countContracts(t, ctx, driverID))

	var status string
	var startDate time.Time
	var id uuid.UUID
	err = testPool.QueryRow(ctx,
		`SELECT id, status, start_date FROM contracts WHERE user_id = $1`,
		driverID,
	).Scan(&id, &status, &startDate)
	require.NoError(t, err)
	assert.Equal(t, existingID, id)
	assert.Equal(t, string(entity.StatusActive), status)
	assert.True(t, startDate.Equal(occurredAt) || startDate.Truncate(time.Second).Equal(occurredAt))
}

func TestProcessor_PublishedCreatesActiveWhenNone(t *testing.T) {
	ctx := context.Background()
	processor := newTestProcessor()
	driverID := uuid.New()
	occurredAt := time.Now().UTC().Truncate(time.Second)

	err := processor.Handle(ctx, TripEvent{
		EventID:    uuid.New().String(),
		EventType:  EventTypeTripPublished,
		DriverID:   driverID.String(),
		CompanyID:  uuid.New().String(),
		OccurredAt: occurredAt,
	})
	require.NoError(t, err)
	require.Equal(t, 1, countContracts(t, ctx, driverID))

	var status string
	err = testPool.QueryRow(ctx,
		`SELECT status FROM contracts WHERE user_id = $1`,
		driverID,
	).Scan(&status)
	require.NoError(t, err)
	assert.Equal(t, string(entity.StatusActive), status)
}

func TestProcessor_PreapprovedDoesNothing(t *testing.T) {
	ctx := context.Background()
	processor := newTestProcessor()
	driverID := uuid.New()

	err := processor.Handle(ctx, TripEvent{
		EventID:   uuid.New().String(),
		EventType: EventType("trip_preapproved"),
		DriverID:  driverID.String(),
		CompanyID: uuid.New().String(),
	})
	require.NoError(t, err)
	require.Equal(t, 0, countContracts(t, ctx, driverID))
}

func TestProcessor_CompletedTerminatesOpenContract(t *testing.T) {
	ctx := context.Background()
	processor := newTestProcessor()
	driverID := uuid.New()
	occurredAt := time.Now().UTC().Truncate(time.Second)
	contractID := uuid.New()

	_, err := testPool.Exec(ctx, `
		INSERT INTO contracts (id, company_id, user_id, status, start_date, created_at)
		VALUES ($1, $2, $3, 'active', now(), now())
	`, contractID, uuid.New(), driverID)
	require.NoError(t, err)

	err = processor.Handle(ctx, TripEvent{
		EventID:    uuid.New().String(),
		EventType:  EventTypeTripCompleted,
		DriverID:   driverID.String(),
		OccurredAt: occurredAt,
	})
	require.NoError(t, err)

	var status string
	var endDate time.Time
	err = testPool.QueryRow(ctx,
		`SELECT status, end_date FROM contracts WHERE id = $1`,
		contractID,
	).Scan(&status, &endDate)
	require.NoError(t, err)
	assert.Equal(t, string(entity.StatusTerminated), status)
	assert.True(t, endDate.Equal(occurredAt) || endDate.Truncate(time.Second).Equal(occurredAt))
}

func TestProcessor_CanceledTerminatesOpenContract(t *testing.T) {
	ctx := context.Background()
	processor := newTestProcessor()
	driverID := uuid.New()
	occurredAt := time.Now().UTC().Truncate(time.Second)
	contractID := uuid.New()

	_, err := testPool.Exec(ctx, `
		INSERT INTO contracts (id, company_id, user_id, status, start_date, created_at)
		VALUES ($1, $2, $3, 'active', now(), now())
	`, contractID, uuid.New(), driverID)
	require.NoError(t, err)

	err = processor.Handle(ctx, TripEvent{
		EventID:    uuid.New().String(),
		EventType:  EventTypeTripCanceled,
		DriverID:   driverID.String(),
		OccurredAt: occurredAt,
	})
	require.NoError(t, err)

	var status string
	var endDate time.Time
	err = testPool.QueryRow(ctx,
		`SELECT status, end_date FROM contracts WHERE id = $1`,
		contractID,
	).Scan(&status, &endDate)
	require.NoError(t, err)
	assert.Equal(t, string(entity.StatusTerminated), status)
	assert.True(t, endDate.Equal(occurredAt) || endDate.Truncate(time.Second).Equal(occurredAt))
}

func TestProcessor_CompletedSkipsWhenNoOpenContract(t *testing.T) {
	ctx := context.Background()
	processor := newTestProcessor()

	err := processor.Handle(ctx, TripEvent{
		EventID:   uuid.New().String(),
		EventType: EventTypeTripCompleted,
		DriverID:  uuid.New().String(),
	})
	require.NoError(t, err)
}

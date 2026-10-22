package events

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"job4j_go_share_trip_contracts/internal/domain/contract/entity"
	"job4j_go_share_trip_contracts/internal/domain/contract/repository"
	"job4j_go_share_trip_contracts/internal/storage"
)

type Processor struct {
	pool            *pgxpool.Pool
	contracts       *repository.UserContractRepository
	processedEvents *repository.ProcessedEventRepository
}

func NewProcessor(
	pool *pgxpool.Pool,
	contracts *repository.UserContractRepository,
	processedEvents *repository.ProcessedEventRepository,
) *Processor {
	return &Processor{
		pool:            pool,
		contracts:       contracts,
		processedEvents: processedEvents,
	}
}

func (p *Processor) Handle(ctx context.Context, event TripEvent) error {
	if event.EventID == "" {
		return fmt.Errorf("missing event_id")
	}
	if event.DriverID == "" {
		return fmt.Errorf("missing driver_id")
	}

	return storage.TxWithoutResult(ctx, p.pool, func(tx pgx.Tx) error {
		inserted, err := p.processedEvents.TryInsert(ctx, tx, event.EventID)
		if err != nil {
			return err
		}
		if !inserted {
			slog.Info("inbox duplicate skipped",
				slog.String("service", "contracts"),
				slog.String("operation", "InboxInsert"),
				slog.String("event_id", event.EventID),
				slog.String("result", "duplicate"),
			)
			return nil
		}

		switch event.EventType {
		case EventTypeTripCreated:
			return p.handleCreated(ctx, tx, event)
		case EventTypeTripPublished:
			return p.handlePublished(ctx, tx, event)
		case EventTypeTripCompleted, EventTypeTripCanceled:
			return p.handleCompleted(ctx, tx, event)
		default:
			return nil
		}
	})
}

func (p *Processor) handleCreated(ctx context.Context, tx pgx.Tx, event TripEvent) error {
	userID, err := uuid.Parse(event.DriverID)
	if err != nil {
		return fmt.Errorf("invalid driver_id: %w", err)
	}
	companyID, err := uuid.Parse(event.CompanyID)
	if err != nil {
		return fmt.Errorf("invalid company_id: %w", err)
	}

	open, err := p.contracts.GetOpenByUserIDForUpdate(ctx, tx, userID)
	if err != nil {
		return err
	}
	if open != nil {
		slog.Info("open driver contract already exists, skip draft create",
			slog.String("contract_id", open.ID.String()),
			slog.String("driver_id", event.DriverID),
			slog.String("event_id", event.EventID),
		)
		return nil
	}

	contract := &entity.Contract{
		ID:        uuid.New(),
		CompanyID: companyID,
		UserID:    userID,
		Status:    entity.StatusDraft,
		CreatedAt: time.Now().UTC(),
	}
	if err := p.contracts.CreateTx(ctx, tx, contract); err != nil {
		return err
	}

	slog.Info("driver contract created as draft",
		slog.String("contract_id", contract.ID.String()),
		slog.String("driver_id", event.DriverID),
		slog.String("event_id", event.EventID),
	)
	return nil
}

func (p *Processor) handlePublished(ctx context.Context, tx pgx.Tx, event TripEvent) error {
	userID, err := uuid.Parse(event.DriverID)
	if err != nil {
		return fmt.Errorf("invalid driver_id: %w", err)
	}
	companyID, err := uuid.Parse(event.CompanyID)
	if err != nil {
		return fmt.Errorf("invalid company_id: %w", err)
	}

	occurredAt := event.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = time.Now().UTC()
	}

	open, err := p.contracts.GetOpenByUserIDForUpdate(ctx, tx, userID)
	if err != nil {
		return err
	}

	if open != nil {
		open.Status = entity.StatusActive
		if open.StartDate == nil {
			open.StartDate = &occurredAt
		}
		if err := p.contracts.UpdateTx(ctx, tx, open); err != nil {
			return err
		}
		slog.Info("existing driver contract activated",
			slog.String("contract_id", open.ID.String()),
			slog.String("driver_id", event.DriverID),
			slog.String("event_id", event.EventID),
		)
		return nil
	}

	contract := &entity.Contract{
		ID:        uuid.New(),
		CompanyID: companyID,
		UserID:    userID,
		Status:    entity.StatusActive,
		StartDate: &occurredAt,
		CreatedAt: time.Now().UTC(),
	}
	if err := p.contracts.CreateTx(ctx, tx, contract); err != nil {
		return err
	}

	slog.Info("driver contract created",
		slog.String("contract_id", contract.ID.String()),
		slog.String("driver_id", event.DriverID),
		slog.String("event_id", event.EventID),
	)
	return nil
}

func (p *Processor) handleCompleted(ctx context.Context, tx pgx.Tx, event TripEvent) error {
	userID, err := uuid.Parse(event.DriverID)
	if err != nil {
		return fmt.Errorf("invalid driver_id: %w", err)
	}

	occurredAt := event.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = time.Now().UTC()
	}

	open, err := p.contracts.GetOpenByUserIDForUpdate(ctx, tx, userID)
	if err != nil {
		return err
	}
	if open == nil {
		slog.Info("no open contract to terminate",
			slog.String("driver_id", event.DriverID),
			slog.String("event_id", event.EventID),
		)
		return nil
	}

	open.Status = entity.StatusTerminated
	open.EndDate = &occurredAt
	if err := p.contracts.UpdateTx(ctx, tx, open); err != nil {
		return err
	}

	slog.Info("driver contract terminated",
		slog.String("contract_id", open.ID.String()),
		slog.String("driver_id", event.DriverID),
		slog.String("event_id", event.EventID),
	)
	return nil
}

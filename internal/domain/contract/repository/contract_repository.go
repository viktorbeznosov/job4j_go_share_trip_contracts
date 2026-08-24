// internal/domain/contract/repository/contract_repository.go
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"job4j_go_share_trip_contracts/internal/domain/contract/entity"
)

type ContractRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Contract, error)
	Create(ctx context.Context, contract *entity.Contract) error
	Update(ctx context.Context, contract *entity.Contract) error
	GetActiveByCompanyID(ctx context.Context, companyID uuid.UUID) (*entity.Contract, error)
}

// MockRepository реализует ContractRepository с моковыми данными
type MockRepository struct {
    db *pgxpool.Pool
	contracts  map[uuid.UUID]*entity.Contract
	companyIDs map[uuid.UUID]uuid.UUID // companyID -> contractID
}

func NewRepository(db *pgxpool.Pool) *MockRepository {
	repo := &MockRepository{
	    db: db,
		contracts:  make(map[uuid.UUID]*entity.Contract),
		companyIDs: make(map[uuid.UUID]uuid.UUID),
	}
	repo.Seed()
	return repo
}

// Seed добавляет тестовые данные
func (r *MockRepository) Seed() {
	companyID1 := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	companyID2 := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	companyID3 := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	companyID4 := uuid.MustParse("44444444-4444-4444-4444-444444444444")

	// Компания 1 - активный контракт со всеми услугами
	contract1 := &entity.Contract{
		ID:        uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
		CompanyID: companyID1,
		Status:    entity.StatusActive,
		ValidFrom: time.Now(),
		ValidTo:   time.Now().AddDate(1, 0, 0),
		Services: []entity.ContractService{
			{Service: entity.ServiceTripCreation, Enabled: true},
			{Service: entity.ServiceTripParticipants, Enabled: true},
			{Service: entity.ServiceNotifications, Enabled: true},
			{Service: entity.ServicePremiumSupport, Enabled: true},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Компания 2 - активный контракт с ограниченными услугами
	contract2 := &entity.Contract{
		ID:        uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
		CompanyID: companyID2,
		Status:    entity.StatusActive,
		ValidFrom: time.Now(),
		ValidTo:   time.Now().AddDate(1, 0, 0),
		Services: []entity.ContractService{
			{Service: entity.ServiceTripCreation, Enabled: true},
			{Service: entity.ServiceTripParticipants, Enabled: false},
			{Service: entity.ServiceNotifications, Enabled: true},
			{Service: entity.ServicePremiumSupport, Enabled: false},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Компания 3 - приостановленный контракт
	contract3 := &entity.Contract{
		ID:        uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc"),
		CompanyID: companyID3,
		Status:    entity.StatusSuspended,
		ValidFrom: time.Now().AddDate(-1, 0, 0),
		ValidTo:   time.Now().AddDate(1, 0, 0),
		Services: []entity.ContractService{
			{Service: entity.ServiceTripCreation, Enabled: true},
			{Service: entity.ServiceNotifications, Enabled: true},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Компания 4 - просроченный контракт
	contract4 := &entity.Contract{
		ID:        uuid.MustParse("dddddddd-dddd-dddd-dddd-dddddddddddd"),
		CompanyID: companyID4,
		Status:    entity.StatusActive,
		ValidFrom: time.Now().AddDate(-2, 0, 0),
		ValidTo:   time.Now().AddDate(-1, 0, 0), // Просрочен
		Services: []entity.ContractService{
			{Service: entity.ServiceTripCreation, Enabled: true},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	r.contracts[contract1.ID] = contract1
	r.contracts[contract2.ID] = contract2
	r.contracts[contract3.ID] = contract3
	r.contracts[contract4.ID] = contract4
	r.companyIDs[companyID1] = contract1.ID
	r.companyIDs[companyID2] = contract2.ID
	// Компания 3 не имеет активного контракта (suspended)
	// Компания 4 не имеет активного контракта (просрочен, но в маппинге не будет)
}

func (r *MockRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Contract, error) {
	contract, exists := r.contracts[id]
	if !exists {
		return nil, fmt.Errorf("contract with id %s not found", id)
	}
	return contract, nil
}

func (r *MockRepository) Create(ctx context.Context, contract *entity.Contract) error {
	if contract.Status == entity.StatusActive {
		if _, exists := r.companyIDs[contract.CompanyID]; exists {
			return fmt.Errorf("company already has an active contract")
		}
	}

	r.contracts[contract.ID] = contract
	
	if contract.Status == entity.StatusActive {
		r.companyIDs[contract.CompanyID] = contract.ID
	}

	return nil
}

func (r *MockRepository) Update(ctx context.Context, contract *entity.Contract) error {
	existing, exists := r.contracts[contract.ID]
	if !exists {
		return fmt.Errorf("contract with id %s not found", contract.ID)
	}

	if contract.Status == entity.StatusActive && existing.Status != entity.StatusActive {
		if activeContractID, exists := r.companyIDs[contract.CompanyID]; exists && activeContractID != contract.ID {
			return fmt.Errorf("company already has an active contract: %s", activeContractID)
		}
	}

	if existing.Status == entity.StatusActive && contract.Status != entity.StatusActive {
		delete(r.companyIDs, contract.CompanyID)
	}

	r.contracts[contract.ID] = contract

	if contract.Status == entity.StatusActive {
		r.companyIDs[contract.CompanyID] = contract.ID
	}

	return nil
}

func (r *MockRepository) GetActiveByCompanyID(ctx context.Context, companyID uuid.UUID) (*entity.Contract, error) {
	contractID, exists := r.companyIDs[companyID]
	if !exists {
		return nil, fmt.Errorf("no active contract found for company %s", companyID)
	}

	contract, exists := r.contracts[contractID]
	if !exists {
		return nil, fmt.Errorf("contract %s not found", contractID)
	}

	if contract.Status != entity.StatusActive {
		delete(r.companyIDs, companyID)
		return nil, fmt.Errorf("contract %s is not active", contractID)
	}

	return contract, nil
}



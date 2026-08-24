package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"job4j_go_share_trip_contracts/internal/domain/contract/entity"
	"job4j_go_share_trip_contracts/internal/domain/contract/repository"
)

type ContractService struct {
	repo repository.ContractRepository
}

func NewContractService(repo repository.ContractRepository) *ContractService {
	return &ContractService{
		repo: repo,
	}
}

func (s *ContractService) GetByID(ctx context.Context, id uuid.UUID) (*entity.Contract, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("contract id is required")
	}

	contract, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract: %w", err)
	}

	return contract, nil
}

func (s *ContractService) Create(ctx context.Context, companyID uuid.UUID, validFrom, validTo time.Time, services []entity.ContractService) (*entity.Contract, error) {
	existing, err := s.repo.GetActiveByCompanyID(ctx, companyID)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("company already has an active contract")
	}

	contract := entity.NewContract(companyID, validFrom, validTo)

	contract.Services = services

	if err := s.repo.Create(ctx, contract); err != nil {
		return nil, fmt.Errorf("failed to create contract: %w", err)
	}

	return contract, nil
}

func (s *ContractService) GetActiveByCompanyID(ctx context.Context, companyID uuid.UUID) (*entity.Contract, error) {
	if companyID == uuid.Nil {
		return nil, fmt.Errorf("company id is required")
	}

	contract, err := s.repo.GetActiveByCompanyID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active contract: %w", err)
	}

	return contract, nil
}

func (s *ContractService) ChangeStatus(ctx context.Context, contractID uuid.UUID, newStatusStr string) (*entity.Contract, error) {
	if contractID == uuid.Nil {
		return nil, fmt.Errorf("contract id is required")
	}

	contract, err := s.repo.GetByID(ctx, contractID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract: %w", err)
	}

	newStatus := entity.ContractStatus(newStatusStr)

	if err := contract.ChangeStatus(newStatus); err != nil {
		return nil, fmt.Errorf("invalid status transition: %w", err)
	}

	if err := s.repo.Update(ctx, contract); err != nil {
		return nil, fmt.Errorf("failed to update contract: %w", err)
	}

	return contract, nil
}

func (s *ContractService) UpdateServices(ctx context.Context, contractID uuid.UUID, services []entity.ContractService) (*entity.Contract, error) {
	if contractID == uuid.Nil {
		return nil, fmt.Errorf("contract id is required")
	}

	if len(services) == 0 {
		return nil, fmt.Errorf("services list cannot be empty")
	}

	contract, err := s.repo.GetByID(ctx, contractID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract: %w", err)
	}

	if contract.Status == entity.StatusTerminated {
		return nil, fmt.Errorf("cannot update services for terminated contract")
	}

	if err := contract.UpdateServices(services); err != nil {
		return nil, fmt.Errorf("failed to update services: %w", err)
	}

	if err := s.repo.Update(ctx, contract); err != nil {
		return nil, fmt.Errorf("failed to update contract: %w", err)
	}

	return contract, nil
}

func (s *ContractService) CheckAvailability(ctx context.Context, companyID uuid.UUID, serviceTypeStr string) (bool, string, error) {
	if companyID == uuid.Nil {
		return false, "", fmt.Errorf("company id is required")
	}

	contract, err := s.repo.GetActiveByCompanyID(ctx, companyID)
	if err != nil {
		return false, "No active contract found for this company", nil
	}

	if contract.IsExpired() {
		return false, "Contract has expired", nil
	}

	if !contract.IsActive() {
		return false, fmt.Sprintf("Contract is not active (current status: %s)", contract.Status), nil
	}

	serviceType := entity.ServiceType(serviceTypeStr)
	enabled := contract.IsServiceEnabled(serviceType)

	if !enabled {
		return false, fmt.Sprintf("Service '%s' is not enabled in the contract", serviceTypeStr), nil
	}

	return true, "", nil
}
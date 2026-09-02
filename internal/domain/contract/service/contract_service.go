package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"

	contractErrors "job4j_go_share_trip_contracts/internal/api/errors"
	"job4j_go_share_trip_contracts/internal/domain/contract/entity"
	"job4j_go_share_trip_contracts/internal/domain/contract/repository"
)

type ContractService struct {
	repo repository.ContractRepository
}

type CreateContractRequest struct {
	CompanyID uuid.UUID
	ValidFrom types.Date
	ValidTo   types.Date
	Services  []ContractServiceInput
}

type ContractServiceInput struct {
	Service string
	Enabled bool
}

type CreateContractResponse struct {
	ID          uuid.UUID
	CompanyID   uuid.UUID
	Status      string
	ValidFrom   types.Date
	ValidTo     types.Date
	Services    []ContractServiceOutput
}

type ContractServiceOutput struct {
	Service string
	Enabled bool
}

type SignContractRequest struct {
	ContractID uuid.UUID
}

type SignContractResponse struct {
	ID          uuid.UUID
	CompanyID   uuid.UUID
	Status      string
	ValidFrom   types.Date
	ValidTo     types.Date
	Services    []ContractServiceOutput
}

type ContractResponse struct {
	ID        uuid.UUID
	CompanyID uuid.UUID
	Status    string
	ValidFrom types.Date
	ValidTo   types.Date
	Services  []ContractServiceOutput
}

type UpdateServicesRequest struct {
	ContractID uuid.UUID
	Services   []ContractServiceInput
}

type UpdateServicesResponse struct {
	ID        uuid.UUID
	CompanyID uuid.UUID
	Status    string
	ValidFrom types.Date
	ValidTo   types.Date
	Services  []ContractServiceOutput
}

func NewContractService(repo repository.ContractRepository) *ContractService {
	return &ContractService{
		repo: repo,
	}
}

func (s *ContractService) GetByID(ctx context.Context, id uuid.UUID) (*ContractResponse, error) {
	if id == uuid.Nil {
		return nil, contractErrors.ErrInvalidContractID
	}

	contract, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toContractResponse(contract), nil
}

func (s *ContractService) Create(ctx context.Context, req CreateContractRequest) (*CreateContractResponse, error) {
	existing, err := s.repo.GetActiveByCompanyID(ctx, req.CompanyID)
	if err == nil && existing != nil {
		return nil, contractErrors.ErrCompanyAlreadyHasActiveContract
	}

	from := req.ValidFrom.Time
	to := req.ValidTo.Time

	if from.After(to) {
		return nil, contractErrors.ErrInvalidDateRange
	}

	services := make([]entity.ContractService, len(req.Services))
	for i, s := range req.Services {
		services[i] = entity.ContractService{
			Service: entity.ServiceType(s.Service),
			Enabled: s.Enabled,
		}
	}

	contract := entity.NewContract(
		req.CompanyID,
		from,
		to,
	)
	contract.Services = services

	if err := s.repo.Create(ctx, contract); err != nil {
		return nil, err
	}

	return s.toCreateContractResponse(contract), nil
}

func (s *ContractService) GetActiveByCompanyID(ctx context.Context, companyID uuid.UUID) (*ContractResponse, error) {
	if companyID == uuid.Nil {
		return nil, contractErrors.ErrInvalidCompanyID
	}

	contract, err := s.repo.GetActiveByCompanyID(ctx, companyID)
	if err != nil {
		return nil, err
	}

	return s.toContractResponse(contract), nil
}

func (s *ContractService) ChangeStatus(ctx context.Context, contractID uuid.UUID, newStatusStr string) (*ContractResponse, error) {
	if contractID == uuid.Nil {
		return nil, contractErrors.ErrInvalidContractID
	}

	contract, err := s.repo.GetByID(ctx, contractID)
	if err != nil {
		return nil, err
	}

	newStatus := entity.ContractStatus(newStatusStr)

	if err := contract.ChangeStatus(newStatus); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, contract); err != nil {
		return nil, err
	}

	return s.toContractResponse(contract), nil
}

func (s *ContractService) UpdateServices(ctx context.Context, req UpdateServicesRequest) (*UpdateServicesResponse, error) {
	if req.ContractID == uuid.Nil {
		return nil, contractErrors.ErrInvalidContractID
	}

	if len(req.Services) == 0 {
		return nil, contractErrors.ErrServicesListEmpty
	}

	contract, err := s.repo.GetByID(ctx, req.ContractID)
	if err != nil {
		return nil, err
	}

	if contract.Status == entity.StatusTerminated {
		return nil, contractErrors.ErrCannotUpdateTerminated
	}

	services := make([]entity.ContractService, len(req.Services))
	for i, s := range req.Services {
		services[i] = entity.ContractService{
			Service: entity.ServiceType(s.Service),
			Enabled: s.Enabled,
		}
	}

	if err := contract.UpdateServices(services); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, contract); err != nil {
		return nil, err
	}

	return s.toUpdateServicesResponse(contract), nil
}

func (s *ContractService) CheckAvailability(ctx context.Context, companyID uuid.UUID, serviceTypeStr string) (bool, string, error) {
	if companyID == uuid.Nil {
		return false, "", contractErrors.ErrInvalidCompanyID
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

func (s *ContractService) SignContract(ctx context.Context, req SignContractRequest) (*SignContractResponse, error) {
	if req.ContractID == uuid.Nil {
		return nil, contractErrors.ErrInvalidContractID
	}

	contract, err := s.repo.GetByID(ctx, req.ContractID)
	if err != nil {
		return nil, err
	}

	if contract.IsExpired() {
		return nil, contractErrors.ErrContractExpired
	}

	if contract.Status != entity.StatusDraft {
		return nil, contractErrors.ErrContractNotDraft
	}

	if err := contract.ChangeStatus(entity.StatusActive); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, contract); err != nil {
		return nil, err
	}

	return s.toSignContractResponse(contract), nil
}

func (s *ContractService) toCreateContractResponse(contract *entity.Contract) *CreateContractResponse {
	services := make([]ContractServiceOutput, len(contract.Services))
	for i, s := range contract.Services {
		services[i] = ContractServiceOutput{
			Service: string(s.Service),
			Enabled: s.Enabled,
		}
	}

	return &CreateContractResponse{
		ID:          contract.ID,
		CompanyID:   contract.CompanyID,
		Status:      string(contract.Status),
		ValidFrom:   types.Date{Time: contract.ValidFrom},
		ValidTo:     types.Date{Time: contract.ValidTo},
		Services:    services,
	}
}

func (s *ContractService) toSignContractResponse(contract *entity.Contract) *SignContractResponse {
	services := make([]ContractServiceOutput, len(contract.Services))
	for i, s := range contract.Services {
		services[i] = ContractServiceOutput{
			Service: string(s.Service),
			Enabled: s.Enabled,
		}
	}

	return &SignContractResponse{
		ID:        contract.ID,
		CompanyID: contract.CompanyID,
		Status:    string(contract.Status),
		ValidFrom: types.Date{Time: contract.ValidFrom},
		ValidTo:   types.Date{Time: contract.ValidTo},
		Services:  services,
	}
}

func (s *ContractService) toContractResponse(contract *entity.Contract) *ContractResponse {
	services := make([]ContractServiceOutput, len(contract.Services))
	for i, s := range contract.Services {
		services[i] = ContractServiceOutput{
			Service: string(s.Service),
			Enabled: s.Enabled,
		}
	}

	return &ContractResponse{
		ID:        contract.ID,
		CompanyID: contract.CompanyID,
		Status:    string(contract.Status),
		ValidFrom: types.Date{Time: contract.ValidFrom},
		ValidTo:   types.Date{Time: contract.ValidTo},
		Services:  services,
	}
}

func (s *ContractService) toUpdateServicesResponse(contract *entity.Contract) *UpdateServicesResponse {
	services := make([]ContractServiceOutput, len(contract.Services))
	for i, s := range contract.Services {
		services[i] = ContractServiceOutput{
			Service: string(s.Service),
			Enabled: s.Enabled,
		}
	}

	return &UpdateServicesResponse{
		ID:        contract.ID,
		CompanyID: contract.CompanyID,
		Status:    string(contract.Status),
		ValidFrom: types.Date{Time: contract.ValidFrom},
		ValidTo:   types.Date{Time: contract.ValidTo},
		Services:  services,
	}
}
package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"job4j_go_share_trip_contracts/gen"
	"job4j_go_share_trip_contracts/internal/domain/contract/service"
)

func (h *ContractHandler) CreateContract(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var req gen.CreateContractRequest
	if err := c.BodyParser(&req); err != nil {
		return h.errorMapper.MapParseError(c, err, "Invalid JSON body")
	}

	if err := validateCreateContractRequest(&req); err != nil {
		return h.errorMapper.MapValidationError(c, err)
	}

	serviceReq := service.CreateContractRequest{
		CompanyID: uuid.UUID(req.CompanyId),
		ValidFrom: req.ValidFrom,
		ValidTo:   req.ValidTo,
		Services:  make([]service.ContractServiceInput, len(*req.Services)),
	}

	for i, s := range *req.Services {
		serviceReq.Services[i] = service.ContractServiceInput{
			Service: string(s.Service),
			Enabled: s.Enabled,
		}
	}

	contractResp, err := h.ContractService.Create(ctx, serviceReq)
	if err != nil {
		return h.errorMapper.MapError(c, err)
	}

	resp := convertToGenCreateContractResponse(contractResp)

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func validateCreateContractRequest(req *gen.CreateContractRequest) error {
	if req.CompanyId == uuid.Nil {
		return fmt.Errorf("companyId is required")
	}

	if req.ValidFrom.IsZero() {
		return fmt.Errorf("validFrom is required")
	}

	if req.ValidTo.IsZero() {
		return fmt.Errorf("validTo is required")
	}

	if req.ValidTo.Before(req.ValidFrom.Time) {
		return fmt.Errorf("validTo must be after or equal to validFrom")
	}

	validServices := map[gen.ServiceType]bool{
		gen.TripCreation:     true,
		gen.TripParticipants: true,
		gen.Notifications:    true,
		gen.PremiumSupport:   true,
	}

	if req.Services != nil {
		for _, s := range *req.Services {
			if !validServices[s.Service] {
				return fmt.Errorf("invalid service type: %s", s.Service)
			}
		}
	}

	return nil
}

func convertToGenCreateContractResponse(resp *service.CreateContractResponse) gen.CreateContractResponse {
	services := make([]gen.ContractService, len(resp.Services))
	for i, s := range resp.Services {
		services[i] = gen.ContractService{
			Service: gen.ServiceType(s.Service),
			Enabled: s.Enabled,
		}
	}

	return gen.CreateContractResponse{
		Id:        resp.ID,
		CompanyId: resp.CompanyID,
		Status:    gen.ContractStatus(resp.Status),
		ValidFrom: resp.ValidFrom,
		ValidTo:   resp.ValidTo,
		Services:  services,
	}
}
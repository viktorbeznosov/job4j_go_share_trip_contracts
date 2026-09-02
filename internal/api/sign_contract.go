package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"job4j_go_share_trip_contracts/gen"
	"job4j_go_share_trip_contracts/internal/domain/contract/service"
)

func (h *ContractHandler) SignContract(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var req gen.SignContractRequest
	if err := c.BodyParser(&req); err != nil {
		return h.errorMapper.MapParseError(c, err, "Invalid JSON body")
	}

	if err := validateSignContractRequest(&req); err != nil {
		return h.errorMapper.MapValidationError(c, err)
	}

	serviceReq := service.SignContractRequest{
		ContractID: uuid.UUID(req.ContractId),
	}

	contractResp, err := h.ContractService.SignContract(ctx, serviceReq)
	if err != nil {
		return h.errorMapper.MapError(c, err)
	}

	resp := convertToGenSignContractResponse(contractResp)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func validateSignContractRequest(req *gen.SignContractRequest) error {
	if req.ContractId == uuid.Nil {
		return fmt.Errorf("contractId is required")
	}
	return nil
}

func convertToGenSignContractResponse(resp *service.SignContractResponse) gen.SignContractResponse {
	services := make([]gen.ContractService, len(resp.Services))
	for i, s := range resp.Services {
		services[i] = gen.ContractService{
			Service: gen.ServiceType(s.Service),
			Enabled: s.Enabled,
		}
	}

	return gen.SignContractResponse{
		Id:        resp.ID,
		CompanyId: resp.CompanyID,
		Status:    gen.ContractStatus(resp.Status),
		ValidFrom: resp.ValidFrom,
		ValidTo:   resp.ValidTo,
		Services:  services,
	}
}
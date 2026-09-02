package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"job4j_go_share_trip_contracts/gen"
	"job4j_go_share_trip_contracts/internal/domain/contract/service"
)

func (h *ContractHandler) GetActiveContract(c *fiber.Ctx, companyId gen.CompanyId) error {
	ctx := c.UserContext()

	companyID := uuid.UUID(companyId)

	if companyID == uuid.Nil {
		return h.errorMapper.MapValidationError(c, fmt.Errorf("companyId is required"))
	}

	contractResp, err := h.ContractService.GetActiveByCompanyID(ctx, companyID)
	if err != nil {
		return h.errorMapper.MapError(c, err)
	}

	resp := convertActiveContractToGenContractResponse(contractResp)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func convertActiveContractToGenContractResponse(resp *service.ContractResponse) gen.ContractResponse {
	services := make([]gen.ContractService, len(resp.Services))
	for i, s := range resp.Services {
		services[i] = gen.ContractService{
			Service: gen.ServiceType(s.Service),
			Enabled: s.Enabled,
		}
	}

	return gen.ContractResponse{
		Id:        resp.ID,
		CompanyId: resp.CompanyID,
		Status:    gen.ContractStatus(resp.Status),
		ValidFrom: resp.ValidFrom,
		ValidTo:   resp.ValidTo,
		Services:  services,
	}
}
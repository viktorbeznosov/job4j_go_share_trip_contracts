package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"job4j_go_share_trip_contracts/gen"
	"job4j_go_share_trip_contracts/internal/domain/contract/service"
)

func (h *ContractHandler) GetContract(c *fiber.Ctx, contractId gen.ContractId) error {
	ctx := c.UserContext()

	id := uuid.UUID(contractId)

	// Получаем DTO от сервиса
	contractResp, err := h.ContractService.GetByID(ctx, id)
	if err != nil {
		return h.errorMapper.MapError(c, err)
	}

	// Конвертируем DTO сервиса в gen.ContractResponse
	resp := convertServiceToGenContractResponse(contractResp)

	return c.Status(fiber.StatusOK).JSON(resp)
}

// convertServiceToGenContractResponse конвертирует service.GetContractResponse в gen.ContractResponse
func convertServiceToGenContractResponse(resp *service.ContractResponse) gen.ContractResponse {
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
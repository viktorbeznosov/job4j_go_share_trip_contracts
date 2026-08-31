// internal/domain/contract/handler/update_services.go
package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"job4j_go_share_trip_contracts/gen"
	"job4j_go_share_trip_contracts/internal/domain/contract/entity"
	"job4j_go_share_trip_contracts/internal/domain/contract/request"
	"job4j_go_share_trip_contracts/internal/domain/contract/response"
)

// UpdateContractServices Добавить или обновить список доступных услуг
// (PUT /contracts/{contractId}/services)
func (h *ContractHandler) UpdateContractServices(c *fiber.Ctx, contractId gen.ContractId) error {
	ctx := c.UserContext()

	// contractId уже является uuid.UUID (через openapi_types.UUID)
	id := uuid.UUID(contractId)

	var req request.UpdateServicesRequest
	if err := c.BodyParser(&req); err != nil {
		return h.errorMapper.MapParseError(c, err, "Invalid JSON body")
	}

	req.ContractID = id

	if err := req.Validate(); err != nil {
		return h.errorMapper.MapValidationError(c, err)
	}

	services := make([]entity.ContractService, len(req.Services))
	for i, s := range req.Services {
		services[i] = entity.ContractService{
			Service: entity.ServiceType(s.Service),
			Enabled: s.Enabled,
		}
	}

	contract, err := h.ContractService.UpdateServices(ctx, id, services)
	if err != nil {
		return h.errorMapper.MapError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response.NewSuccessResponse(
		convertToServiceResponse(contract.Services),
	))
}

// convertToServiceResponse конвертирует доменные услуги в ответ
func convertToServiceResponse(services []entity.ContractService) []response.ContractServiceResponse {
	result := make([]response.ContractServiceResponse, len(services))
	for i, s := range services {
		result[i] = response.ContractServiceResponse{
			Service: string(s.Service),
			Enabled: s.Enabled,
		}
	}
	return result
}
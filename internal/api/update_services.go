package api

import (

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"job4j_go_share_trip_contracts/internal/domain/contract/entity"
	"job4j_go_share_trip_contracts/internal/domain/contract/request"
	"job4j_go_share_trip_contracts/internal/domain/contract/response"
)

func (h *ContractHandler) UpdateServices(c *fiber.Ctx) error {
	ctx := c.UserContext()

	contractIDStr := c.Params("contractId")

	contractID, err := uuid.Parse(contractIDStr)
	if err != nil {
		return h.errorMapper.MapParseError(c, err, "Invalid contract ID format")
	}

	var req request.UpdateServicesRequest
	if err := c.BodyParser(&req); err != nil {
		return h.errorMapper.MapParseError(c, err, "Invalid JSON body")
	}

	req.ContractID = contractID

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

	contract, err := h.ContractService.UpdateServices(ctx, contractID, services)
	if err != nil {
		return h.errorMapper.MapError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response.NewSuccessResponse(
		convertToServiceResponse(contract.Services),
	))
}

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



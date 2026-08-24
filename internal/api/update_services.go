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
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrorResponse(
			"Invalid contract ID format",
			err.Error(),
		))
	}

	var req request.UpdateServicesRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrorResponse(
			"Invalid JSON body",
			err.Error(),
		))
	}

	req.ContractID = contractID

	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrorResponse(
			err.Error(),
		))
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
		if err.Error() == "contract with id "+contractID.String()+" not found" {
			return c.Status(fiber.StatusNotFound).JSON(response.NewErrorResponse(
				"Contract not found",
				err.Error(),
			))
		}

		if err.Error() == "cannot update services for terminated contract" ||
			err.Error() == "services list cannot be empty" {
			return c.Status(fiber.StatusConflict).JSON(response.NewErrorResponse(
				"Invalid operation",
				err.Error(),
			))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrorResponse(
			"Failed to update services",
			err.Error(),
		))
	}

	// 7. Возвращаем ответ (список услуг)
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



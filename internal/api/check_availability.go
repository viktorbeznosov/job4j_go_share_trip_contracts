package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"job4j_go_share_trip_contracts/internal/domain/contract/request"
	"job4j_go_share_trip_contracts/internal/domain/contract/response"
)

func (h *ContractHandler) CheckAvailability(c *fiber.Ctx) error {
	ctx := c.UserContext()

	companyIDStr := c.Params("companyId")
	service := c.Params("service")

	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrorResponse(
			"Invalid company ID format",
			err.Error(),
		))
	}

	req := &request.CheckAvailabilityRequest{
		CompanyID: companyID,
		Service:   service,
	}
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrorResponse(
			err.Error(),
		))
	}

	available, reason, err := h.ContractService.CheckAvailability(ctx, companyID, service)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrorResponse(
			"Failed to check availability",
			err.Error(),
		))
	}

	return c.Status(fiber.StatusOK).JSON(response.NewAvailabilityResponse(available, reason))
}




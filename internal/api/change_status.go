package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"job4j_go_share_trip_contracts/internal/domain/contract/request"
	"job4j_go_share_trip_contracts/internal/domain/contract/response"
)

func (h *ContractHandler) ChangeStatus(c *fiber.Ctx) error {
	ctx := c.UserContext()

	contractIDStr := c.Params("contractId")

	contractID, err := uuid.Parse(contractIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrorResponse(
			"Invalid contract ID format",
			err.Error(),
		))
	}

	var req request.ChangeStatusRequest
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

	contract, err := h.ContractService.ChangeStatus(ctx, contractID, req.Status)
	if err != nil {
		if err.Error() == "contract with id "+contractID.String()+" not found" {
			return c.Status(fiber.StatusNotFound).JSON(response.NewErrorResponse(
				"Contract not found",
				err.Error(),
			))
		}

		if err.Error() == "cannot change status from terminated to "+req.Status ||
			err.Error() == "invalid status transition" {
			return c.Status(fiber.StatusConflict).JSON(response.NewErrorResponse(
				"Invalid status transition",
				err.Error(),
			))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrorResponse(
			"Failed to change status",
			err.Error(),
		))
	}

	return c.Status(fiber.StatusOK).JSON(response.NewSuccessResponse(
		response.FromEntity(contract),
	))
}



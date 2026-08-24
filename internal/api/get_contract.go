package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"job4j_go_share_trip_contracts/internal/domain/contract/request"
	"job4j_go_share_trip_contracts/internal/domain/contract/response"
)

func (h *ContractHandler) GetContract(c *fiber.Ctx) error {
	ctx := c.UserContext()

	contractIDStr := c.Params("contractId")

	contractID, err := uuid.Parse(contractIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrorResponse(
			"Invalid contract ID format",
			err.Error(),
		))
	}

	req := &request.GetContractRequest{
		ContractID: contractID,
	}
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrorResponse(
			err.Error(),
		))
	}

	contract, err := h.ContractService.GetByID(ctx, contractID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.NewErrorResponse(
			"Contract not found",
			err.Error(),
		))
	}

	return c.Status(fiber.StatusOK).JSON(response.NewSuccessResponse(
		response.FromEntity(contract),
	))
}
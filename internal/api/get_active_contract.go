package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"job4j_go_share_trip_contracts/internal/domain/contract/request"
	"job4j_go_share_trip_contracts/internal/domain/contract/response"
)

func (h *ContractHandler) GetActiveContract(c *fiber.Ctx) error {
	ctx := c.UserContext()

	companyIDStr := c.Params("companyId")

	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrorResponse(
			"Invalid company ID format",
			err.Error(),
		))
	}

	req := &request.GetActiveContractRequest{
		CompanyID: companyID,
	}
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrorResponse(
			err.Error(),
		))
	}

	contract, err := h.ContractService.GetActiveByCompanyID(ctx, companyID)
	if err != nil {
		if err.Error() == "no active contract found for company "+companyID.String() ||
			err.Error() == "contract "+contract.ID.String()+" is not active" {
			return c.Status(fiber.StatusNotFound).JSON(response.NewErrorResponse(
				"Active contract not found",
				err.Error(),
			))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrorResponse(
			"Failed to get active contract",
			err.Error(),
		))
	}

	// 5. Возвращаем ответ
	return c.Status(fiber.StatusOK).JSON(response.NewSuccessResponse(
		response.FromEntity(contract),
	))
}



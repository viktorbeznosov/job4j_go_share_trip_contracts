package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"job4j_go_share_trip_contracts/gen"
	"job4j_go_share_trip_contracts/internal/domain/contract/response"
)

// GetContract Получить договор по идентификатору
// (GET /contracts/{contractId})
func (h *ContractHandler) GetContract(c *fiber.Ctx, contractId gen.ContractId) error {
	ctx := c.UserContext()

	// contractId уже является uuid.UUID (через openapi_types.UUID)
	id := uuid.UUID(contractId)

	contract, err := h.ContractService.GetByID(ctx, id)
	if err != nil {
		return h.errorMapper.MapError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response.NewSuccessResponse(
		response.FromEntity(contract),
	))
}
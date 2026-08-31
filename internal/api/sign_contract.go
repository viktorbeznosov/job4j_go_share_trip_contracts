package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"job4j_go_share_trip_contracts/internal/domain/contract/request"
	"job4j_go_share_trip_contracts/internal/domain/contract/response"
)

// SignContract Подписание контракта
// (POST /contracts/sign_contract)
func (h *ContractHandler) SignContract(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var req request.SignContractRequest
	if err := c.BodyParser(&req); err != nil {
		return h.errorMapper.MapParseError(c, err, "Invalid JSON body")
	}

	if err := req.Validate(); err != nil {
		return h.errorMapper.MapValidationError(c, err)
	}

	// Парсим UUID
	contractID, err := uuid.Parse(req.ContractID)
	if err != nil {
		return h.errorMapper.MapParseError(c, err, "Invalid contract ID format")
	}

	// Подписываем контракт
	contract, err := h.ContractService.SignContract(ctx, contractID)
	if err != nil {
		return h.errorMapper.MapError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response.NewSuccessResponse(
		response.NewSignContractResponse(contract),
	))
}



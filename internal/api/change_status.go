package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"job4j_go_share_trip_contracts/gen"
	"job4j_go_share_trip_contracts/internal/domain/contract/request"
	"job4j_go_share_trip_contracts/internal/domain/contract/response"
)

// ChangeContractStatus Изменить статус договора
// (PATCH /contracts/{contractId}/status)
func (h *ContractHandler) ChangeContractStatus(c *fiber.Ctx, contractId gen.ContractId) error {
	ctx := c.UserContext()

	// contractId уже является uuid.UUID (через openapi_types.UUID)
	id := uuid.UUID(contractId)

	var req request.ChangeStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return h.errorMapper.MapParseError(c, err, "Invalid JSON body")
	}

	req.ContractID = id

	if err := req.Validate(); err != nil {
		return h.errorMapper.MapValidationError(c, err)
	}

	contract, err := h.ContractService.ChangeStatus(ctx, id, req.Status)
	if err != nil {
		return h.errorMapper.MapError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response.NewSuccessResponse(
		response.FromEntity(contract),
	))
}
package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"job4j_go_share_trip_contracts/gen"
	"job4j_go_share_trip_contracts/internal/domain/contract/request"
	"job4j_go_share_trip_contracts/internal/domain/contract/response"
)

func (h *ContractHandler) GetActiveContract(c *fiber.Ctx, companyId gen.CompanyId) error {
	ctx := c.UserContext()

	companyID := uuid.UUID(companyId)

	req := &request.GetActiveContractRequest{
		CompanyID: companyID,
	}
	if err := req.Validate(); err != nil {
		return h.errorMapper.MapValidationError(c, err)
	}

	contract, err := h.ContractService.GetActiveByCompanyID(ctx, companyID)
	if err != nil {
		return h.errorMapper.MapError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response.NewSuccessResponse(
		response.FromEntity(contract),
	))
}
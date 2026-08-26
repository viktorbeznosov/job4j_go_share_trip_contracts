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
		return h.errorMapper.MapParseError(c, err, "Invalid company ID format")
	}

	req := &request.CheckAvailabilityRequest{
		CompanyID: companyID,
		Service:   service,
	}
	if err := req.Validate(); err != nil {
		return h.errorMapper.MapValidationError(c, err)
	}

	available, reason, err := h.ContractService.CheckAvailability(ctx, companyID, service)
	if err != nil {
		return h.errorMapper.MapError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response.NewAvailabilityResponse(available, reason))
}




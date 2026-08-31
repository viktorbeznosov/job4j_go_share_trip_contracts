package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"job4j_go_share_trip_contracts/gen"
	"job4j_go_share_trip_contracts/internal/domain/contract/request"
	"job4j_go_share_trip_contracts/internal/domain/contract/response"
)

// CheckServiceAvailability Проверить доступность услуги для компании
// (GET /companies/{companyId}/services/{service}/availability)
func (h *ContractHandler) CheckServiceAvailability(c *fiber.Ctx, companyId gen.CompanyId, service gen.Service) error {
	ctx := c.UserContext()

	// companyId уже является uuid.UUID (через openapi_types.UUID)
	companyID := uuid.UUID(companyId)

	// service уже является валидным ServiceType
	serviceType := string(service)

	req := &request.CheckAvailabilityRequest{
		CompanyID: companyID,
		Service:   serviceType,
	}
	if err := req.Validate(); err != nil {
		return h.errorMapper.MapValidationError(c, err)
	}

	available, reason, err := h.ContractService.CheckAvailability(ctx, companyID, serviceType)
	if err != nil {
		return h.errorMapper.MapError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response.NewAvailabilityResponse(available, reason))
}
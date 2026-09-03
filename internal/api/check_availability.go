package api

import (
    "fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"job4j_go_share_trip_contracts/gen"
)

func (h *ContractHandler) CheckServiceAvailability(c *fiber.Ctx, companyId gen.CompanyId, service gen.Service) error {
	ctx := c.UserContext()

	companyID := uuid.UUID(companyId)
	serviceType := string(service)

	if err := validateCheckAvailability(companyID, serviceType); err != nil {
		return h.errorMapper.MapValidationError(c, err)
	}

	available, reason, err := h.ContractService.CheckAvailability(ctx, companyID, serviceType)
	if err != nil {
		return h.errorMapper.MapError(c, err)
	}

	resp := gen.AvailabilityResponse{
		Available: available,
		Reason:    &reason, // reason может быть пустой строкой
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func validateCheckAvailability(companyID uuid.UUID, serviceType string) error {
	if companyID == uuid.Nil {
		return fmt.Errorf("companyId is required")
	}

	if serviceType == "" {
		return fmt.Errorf("service is required")
	}

	validServices := map[string]bool{
		"trip_creation":     true,
		"trip_participants": true,
		"notifications":     true,
		"premium_support":   true,
	}

	if !validServices[serviceType] {
		return fmt.Errorf("invalid service: %s. Must be one of: trip_creation, trip_participants, notifications, premium_support", serviceType)
	}

	return nil
}
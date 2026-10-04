package api

import (
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"job4j_go_share_trip_contracts/gen"
)

func (h *ContractHandler) CheckServiceAvailability(c *fiber.Ctx, companyId gen.CompanyId, service gen.Service) error {
	ctx := c.UserContext()

	companyID := uuid.UUID(companyId)
	serviceType := string(service)
	requestID := firstHeader(c, "X-Request-ID", "X-Request-Id")
	correlationID := firstHeader(c, "X-Correlation-ID")
	if correlationID == "" {
		correlationID = requestID
	}
	tripID := firstHeader(c, "X-Trip-ID")
	userID := firstHeader(c, "X-User-ID")

	if err := validateCheckAvailability(companyID, service); err != nil {
		slog.Info("permission check rejected",
			slog.String("service", "contract"),
			slog.String("operation", "CheckPermission"),
			slog.String("request_id", requestID),
			slog.String("correlation_id", correlationID),
			slog.String("trip_id", tripID),
			slog.String("user_id", userID),
			slog.String("result", "error"),
			slog.String("error", err.Error()),
		)
		return h.errorMapper.MapValidationError(c, err)
	}

	available, reason, err := h.ContractService.CheckAvailability(ctx, companyID, serviceType)
	if err != nil {
		slog.Info("permission check failed",
			slog.String("service", "contract"),
			slog.String("operation", "CheckPermission"),
			slog.String("request_id", requestID),
			slog.String("correlation_id", correlationID),
			slog.String("trip_id", tripID),
			slog.String("user_id", userID),
			slog.String("result", "error"),
			slog.String("error", err.Error()),
		)
		return h.errorMapper.MapError(c, err)
	}

	result := "denied"
	if available {
		result = "allowed"
	}

	slog.Info("permission check completed",
		slog.String("service", "contract"),
		slog.String("operation", "CheckPermission"),
		slog.String("request_id", requestID),
		slog.String("correlation_id", correlationID),
		slog.String("trip_id", tripID),
		slog.String("user_id", userID),
		slog.String("result", result),
	)

	resp := gen.AvailabilityResponse{
		Available: available,
		Reason:    &reason,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func firstHeader(c *fiber.Ctx, names ...string) string {
	for _, name := range names {
		if value := c.Get(name); value != "" {
			return value
		}
	}
	return ""
}

func validateCheckAvailability(companyID uuid.UUID, serviceType gen.ServiceType) error {
	if companyID == uuid.Nil {
		return fmt.Errorf("companyId is required")
	}

	if serviceType == "" {
		return fmt.Errorf("service is required")
	}

	if !serviceType.Valid() {
		return fmt.Errorf("invalid service: %s. Must be one of: trip_creation, trip_participants, notifications, premium_support", serviceType)
	}

	return nil
}

package request

import (
	"fmt"

	"github.com/google/uuid"
)

type CheckAvailabilityRequest struct {
	CompanyID uuid.UUID `params:"companyId"`
	Service   string    `params:"service"`
}

func (r *CheckAvailabilityRequest) Validate() error {
	if r.CompanyID == uuid.Nil {
		return fmt.Errorf("companyId is required")
	}

	if r.Service == "" {
		return fmt.Errorf("service is required")
	}

	validServices := map[string]bool{
		"trip_creation":     true,
		"trip_participants": true,
		"notifications":     true,
		"premium_support":   true,
	}

	if !validServices[r.Service] {
		return fmt.Errorf("invalid service: %s. Must be one of: trip_creation, trip_participants, notifications, premium_support", r.Service)
	}

	return nil
}
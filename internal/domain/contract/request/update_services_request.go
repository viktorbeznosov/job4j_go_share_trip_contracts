package request

import (
	"fmt"

	"github.com/google/uuid"
)

type ServiceInput struct {
	Service string `json:"service"`
	Enabled bool   `json:"enabled"`
}

type UpdateServicesRequest struct {
	ContractID uuid.UUID     `params:"contractId"`
	Services   []ServiceInput `json:"services"`
}

func (r *UpdateServicesRequest) Validate() error {
	if r.ContractID == uuid.Nil {
		return fmt.Errorf("contractId is required")
	}

	if len(r.Services) == 0 {
		return fmt.Errorf("services list cannot be empty")
	}

	validServices := map[string]bool{
		"trip_creation":     true,
		"trip_participants": true,
		"notifications":     true,
		"premium_support":   true,
	}

	seen := make(map[string]bool)
	for _, s := range r.Services {
		if s.Service == "" {
			return fmt.Errorf("service name is required")
		}

		if !validServices[s.Service] {
			return fmt.Errorf("invalid service type: %s. Must be one of: trip_creation, trip_participants, notifications, premium_support", s.Service)
		}

		// Проверяем дубликаты
		if seen[s.Service] {
			return fmt.Errorf("duplicate service: %s", s.Service)
		}
		seen[s.Service] = true
	}

	return nil
}
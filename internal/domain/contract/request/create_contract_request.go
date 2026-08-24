package request

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ContractServiceInput struct {
	Service string `json:"service"`
	Enabled bool   `json:"enabled"`
}

type CreateContractRequest struct {
	CompanyID   string                  `json:"companyId"`
	ValidFrom   string                  `json:"validFrom"`
	ValidTo     string                  `json:"validTo"`
	Services    []ContractServiceInput  `json:"services,omitempty"`
}

func (r *CreateContractRequest) Validate() error {
	if r.CompanyID == "" {
		return fmt.Errorf("companyId is required")
	}
	
	_, err := uuid.Parse(r.CompanyID)
	if err != nil {
		return fmt.Errorf("companyId must be a valid UUID")
	}

	if r.ValidFrom == "" {
		return fmt.Errorf("validFrom is required")
	}
	
	from, err := time.Parse("2006-01-02", r.ValidFrom)
	if err != nil {
		return fmt.Errorf("validFrom must be in format YYYY-MM-DD")
	}

	if r.ValidTo == "" {
		return fmt.Errorf("validTo is required")
	}

	to, err := time.Parse("2006-01-02", r.ValidTo)
	if err != nil {
		return fmt.Errorf("validTo must be in format YYYY-MM-DD")
	}

	if to.Before(from) {
		return fmt.Errorf("validTo must be after or equal to validFrom")
	}

	for _, s := range r.Services {
		if s.Service == "" {
			return fmt.Errorf("service name is required")
		}
		validServices := map[string]bool{
			"trip_creation":     true,
			"trip_participants": true,
			"notifications":     true,
			"premium_support":   true,
		}
		if !validServices[s.Service] {
			return fmt.Errorf("invalid service type: %s", s.Service)
		}
	}

	return nil
}

func (r *CreateContractRequest) ParseDates() (time.Time, time.Time, error) {
	from, err := time.Parse("2006-01-02", r.ValidFrom)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	
	to, err := time.Parse("2006-01-02", r.ValidTo)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	
	return from, to, nil
}
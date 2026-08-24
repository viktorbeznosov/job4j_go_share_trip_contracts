package request

import (
	"fmt"

	"github.com/google/uuid"
)

type ChangeStatusRequest struct {
	ContractID uuid.UUID `params:"contractId"`
	Status     string    `json:"status"`
}

func (r *ChangeStatusRequest) Validate() error {
	if r.ContractID == uuid.Nil {
		return fmt.Errorf("contractId is required")
	}

	if r.Status == "" {
		return fmt.Errorf("status is required")
	}

	validStatuses := map[string]bool{
		"draft":     true,
		"active":    true,
		"suspended": true,
		"terminated": true,
	}

	if !validStatuses[r.Status] {
		return fmt.Errorf("invalid status: %s. Must be one of: draft, active, suspended, terminated", r.Status)
	}

	return nil
}
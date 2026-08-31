package request

import (
	"fmt"

	"github.com/google/uuid"
)

type SignContractRequest struct {
	ContractID string `json:"contractId" validate:"required"`
}

func (r *SignContractRequest) Validate() error {
	if r.ContractID == "" {
		return fmt.Errorf("contractId is required")
	}

	if _, err := uuid.Parse(r.ContractID); err != nil {
		return fmt.Errorf("contractId must be a valid UUID")
	}

	return nil
}
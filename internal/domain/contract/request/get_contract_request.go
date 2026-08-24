package request

import (
	"fmt"

	"github.com/google/uuid"
)

type GetContractRequest struct {
	ContractID uuid.UUID `params:"contractId"`
}

func (r *GetContractRequest) Validate() error {
	if r.ContractID == uuid.Nil {
		return fmt.Errorf("contractId is required")
	}
	return nil
}
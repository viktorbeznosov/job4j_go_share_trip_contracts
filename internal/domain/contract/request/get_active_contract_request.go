package request

import (
	"fmt"

	"github.com/google/uuid"
)

type GetActiveContractRequest struct {
	CompanyID uuid.UUID `params:"companyId"`
}

func (r *GetActiveContractRequest) Validate() error {
	if r.CompanyID == uuid.Nil {
		return fmt.Errorf("companyId is required")
	}
	return nil
}
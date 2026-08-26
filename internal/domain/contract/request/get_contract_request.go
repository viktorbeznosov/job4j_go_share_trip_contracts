package request

import (
	"github.com/google/uuid"
)

type GetContractRequest struct {
	ContractID uuid.UUID `params:"contractId"`
}

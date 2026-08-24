package response

import (
	"job4j_go_share_trip_contracts/internal/domain/contract/entity"
)

type CreateContractResponse struct {
	ID          string                    `json:"id"`
	CompanyID   string                    `json:"companyId"`
	Status      string                    `json:"status"`
	ValidFrom   string                    `json:"validFrom"`
	ValidTo     string                    `json:"validTo"`
	Services    []ContractServiceResponse `json:"services"`
}

func NewCreateContractResponse(contract *entity.Contract) CreateContractResponse {
	services := make([]ContractServiceResponse, len(contract.Services))
	for i, s := range contract.Services {
		services[i] = ContractServiceResponse{
			Service: string(s.Service),
			Enabled: s.Enabled,
		}
	}

	return CreateContractResponse{
		ID:        contract.ID.String(),
		CompanyID: contract.CompanyID.String(),
		Status:    string(contract.Status),
		ValidFrom: contract.ValidFrom.Format("2006-01-02"),
		ValidTo:   contract.ValidTo.Format("2006-01-02"),
		Services:  services,
	}
}
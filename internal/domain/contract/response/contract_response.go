package response

import (
	"job4j_go_share_trip_contracts/internal/domain/contract/entity"
)

type ContractServiceResponse struct {
	Service string `json:"service"`
	Enabled bool   `json:"enabled"`
}

type ContractResponse struct {
	ID          string                    `json:"id"`
	CompanyID   string                    `json:"companyId"`
	Status      string                    `json:"status"`
	ValidFrom   string                    `json:"validFrom"`
	ValidTo     string                    `json:"validTo"`
	Services    []ContractServiceResponse `json:"services"`
}

// FromEntity конвертирует сущность в ответ
func FromEntity(contract *entity.Contract) ContractResponse {
	services := make([]ContractServiceResponse, len(contract.Services))
	for i, s := range contract.Services {
		services[i] = ContractServiceResponse{
			Service: string(s.Service),
			Enabled: s.Enabled,
		}
	}

	return ContractResponse{
		ID:        contract.ID.String(),
		CompanyID: contract.CompanyID.String(),
		Status:    string(contract.Status),
		ValidFrom: contract.ValidFrom.Format("2006-01-02"),
		ValidTo:   contract.ValidTo.Format("2006-01-02"),
		Services:  services,
	}
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

func NewSuccessResponse(data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Success: true,
		Data:    data,
	}
}

func NewErrorResponse(err string, details ...string) *ErrorResponse {
	resp := &ErrorResponse{
		Success: false,
		Error:   err,
	}
	if len(details) > 0 && details[0] != "" {
		resp.Details = details[0]
	}
	return resp
}
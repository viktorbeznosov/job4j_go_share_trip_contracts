package api

import (
	"job4j_go_share_trip_contracts/internal/domain/contract/service"
)

type ContractHandler struct {
	ContractService *service.ContractService
}

func NewContractHandler(contractService *service.ContractService) *ContractHandler {
	return &ContractHandler{
		ContractService: contractService,
	}
}
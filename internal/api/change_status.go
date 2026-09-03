package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"job4j_go_share_trip_contracts/gen"
	contractErrors "job4j_go_share_trip_contracts/internal/api/errors"
	"job4j_go_share_trip_contracts/internal/domain/contract/service"
)

func (h *ContractHandler) ChangeContractStatus(c *fiber.Ctx, contractId gen.ContractId) error {
	ctx := c.UserContext()

	id := uuid.UUID(contractId)

	var req gen.ChangeStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return h.errorMapper.MapParseError(c, err, "Invalid JSON body")
	}


	if err := validateStatus(req.Status); err != nil {
		return h.errorMapper.MapError(c, err)
	}

	contractResp, err := h.ContractService.ChangeStatus(ctx, id, string(req.Status))
	if err != nil {
		return h.errorMapper.MapError(c, err)
	}

	resp := convertChangeStatusToGenContractResponse(contractResp)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func convertChangeStatusToGenContractResponse(resp *service.ContractResponse) gen.ContractResponse {
	services := make([]gen.ContractService, len(resp.Services))
	for i, s := range resp.Services {
		services[i] = gen.ContractService{
			Service: gen.ServiceType(s.Service),
			Enabled: s.Enabled,
		}
	}

	return gen.ContractResponse{
		Id:        resp.ID,
		CompanyId: resp.CompanyID,
		Status:    gen.ContractStatus(resp.Status),
		ValidFrom: resp.ValidFrom,
		ValidTo:   resp.ValidTo,
		Services:  services,
	}
}

func validateStatus(status gen.ContractStatus) error {
	if !status.Valid() {
		return contractErrors.ErrUnknownStatus
	}
	return nil
}
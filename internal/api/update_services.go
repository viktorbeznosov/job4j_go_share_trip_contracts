// internal/domain/contract/handler/update_services.go
package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"job4j_go_share_trip_contracts/gen"
	"job4j_go_share_trip_contracts/internal/domain/contract/service"
)

func (h *ContractHandler) UpdateContractServices(c *fiber.Ctx, contractId gen.ContractId) error {
	ctx := c.UserContext()

	id := uuid.UUID(contractId)

	var req gen.UpdateServicesRequest
	if err := c.BodyParser(&req); err != nil {
		return h.errorMapper.MapParseError(c, err, "Invalid JSON body")
	}

	serviceReq := service.UpdateServicesRequest{
		ContractID: id,
		Services:   make([]service.ContractServiceInput, len(req.Services)),
	}

	for i, s := range req.Services {
		serviceReq.Services[i] = service.ContractServiceInput{
			Service: string(s.Service),
			Enabled: s.Enabled,
		}
	}

	contractResp, err := h.ContractService.UpdateServices(ctx, serviceReq)
	if err != nil {
		return h.errorMapper.MapError(c, err)
	}

	resp := convertUpdateServicesToGenContractResponse(contractResp)

	return c.Status(fiber.StatusOK).JSON(resp)
}

func convertUpdateServicesToGenContractResponse(resp *service.UpdateServicesResponse) gen.ContractResponse {
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
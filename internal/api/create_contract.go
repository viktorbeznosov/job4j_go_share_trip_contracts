package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"job4j_go_share_trip_contracts/internal/domain/contract/entity"
	"job4j_go_share_trip_contracts/internal/domain/contract/request"
	"job4j_go_share_trip_contracts/internal/domain/contract/response"
)

func (h *ContractHandler) CreateContract(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var req request.CreateContractRequest
	if err := c.BodyParser(&req); err != nil {
		return h.errorMapper.MapParseError(c, err, "Invalid JSON body")
	}

	if err := req.Validate(); err != nil {
		return h.errorMapper.MapValidationError(c, err)
	}

	companyID, err := uuid.Parse(req.CompanyID)
	if err != nil {
		return h.errorMapper.MapParseError(c, err, "Invalid company ID format")
	}

	validFrom, validTo, err := req.ParseDates()
	if err != nil {
		return h.errorMapper.MapValidationError(c, err)
	}

	services := make([]entity.ContractService, len(req.Services))
	for i, s := range req.Services {
		services[i] = entity.ContractService{
			Service: entity.ServiceType(s.Service),
			Enabled: s.Enabled,
		}
	}

	contract, err := h.ContractService.Create(ctx, companyID, validFrom, validTo, services)
	if err != nil {
		return h.errorMapper.MapError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(response.NewSuccessResponse(
		response.NewCreateContractResponse(contract),
	))
}
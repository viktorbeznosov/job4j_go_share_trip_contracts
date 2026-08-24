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

	// 1. Парсим JSON
	var req request.CreateContractRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrorResponse(
			"Invalid JSON body",
			err.Error(),
		))
	}

	// 2. Валидируем запрос
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrorResponse(
			err.Error(),
		))
	}

	// 3. Парсим UUID компании
	companyID, err := uuid.Parse(req.CompanyID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrorResponse(
			"Invalid company ID format",
			err.Error(),
		))
	}

	// 4. Парсим даты
	validFrom, validTo, err := req.ParseDates()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.NewErrorResponse(
			err.Error(),
		))
	}

	// 5. Конвертируем услуги
	services := make([]entity.ContractService, len(req.Services))
	for i, s := range req.Services {
		services[i] = entity.ContractService{
			Service: entity.ServiceType(s.Service),
			Enabled: s.Enabled,
		}
	}

	// 6. Создаем контракт через сервис
	contract, err := h.ContractService.Create(ctx, companyID, validFrom, validTo, services)
	if err != nil {

		// Проверяем, не конфликт ли это
		if err.Error() == "company already has an active contract" {
			return c.Status(fiber.StatusConflict).JSON(response.NewErrorResponse(
				"Company already has an active contract",
				err.Error(),
			))
		}
		
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrorResponse(
			"Failed to create contract",
			err.Error(),
		))
	}

	// 7. Возвращаем ответ
	return c.Status(fiber.StatusCreated).JSON(response.NewSuccessResponse(
		response.NewCreateContractResponse(contract),
	))
}
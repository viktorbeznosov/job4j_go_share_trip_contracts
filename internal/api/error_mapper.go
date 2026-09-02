package api

import (
	"github.com/gofiber/fiber/v2"

	"job4j_go_share_trip_contracts/internal/domain/contract/response"
	contractErrors "job4j_go_share_trip_contracts/internal/api/errors"
)

type ErrorMapper struct{}

func NewErrorMapper() *ErrorMapper {
	return &ErrorMapper{}
}

func (m *ErrorMapper) MapError(c *fiber.Ctx, err error) error {
	if err == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.NewErrorResponse(
			"Unknown error",
		))
	}

	status := contractErrors.GetHTTPStatus(err)
	message := contractErrors.GetErrorMessage(err)

	return c.Status(status).JSON(response.NewErrorResponse(
		message,
		err.Error(),
	))
}

func (m *ErrorMapper) MapValidationError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusBadRequest).JSON(response.NewErrorResponse(
		"Validation error",
		err.Error(),
	))
}

func (m *ErrorMapper) MapParseError(c *fiber.Ctx, err error, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(response.NewErrorResponse(
		message,
		err.Error(),
	))
}
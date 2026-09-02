package api

import (
	"github.com/gofiber/fiber/v2"

	contractErrors "job4j_go_share_trip_contracts/internal/api/errors"
)

type ErrorMapper struct{}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

func NewErrorMapper() *ErrorMapper {
	return &ErrorMapper{}
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

func (m *ErrorMapper) MapError(c *fiber.Ctx, err error) error {
	if err == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(NewErrorResponse(
			"Unknown error",
		))
	}

	status := contractErrors.GetHTTPStatus(err)
	message := contractErrors.GetErrorMessage(err)

	return c.Status(status).JSON(NewErrorResponse(
		message,
		err.Error(),
	))
}

func (m *ErrorMapper) MapValidationError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusBadRequest).JSON(NewErrorResponse(
		"Validation error",
		err.Error(),
	))
}

func (m *ErrorMapper) MapParseError(c *fiber.Ctx, err error, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(NewErrorResponse(
		message,
		err.Error(),
	))
}
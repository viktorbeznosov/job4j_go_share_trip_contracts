package api

import (
	"github.com/gofiber/fiber/v2"
)

func (s *Server) Route(route fiber.Router) {
	route.Get("/ready", s.Ready)

    contracts := route.Group("/contracts")
    contracts.Post("/", s.ContractHandler.CreateContract)
    contracts.Get("/:contractId", s.ContractHandler.GetContract)
    contracts.Patch("/:contractId/status", s.ContractHandler.ChangeStatus)
    contracts.Put("/:contractId/services", s.ContractHandler.UpdateServices)

    companies := route.Group("/companies")
    companies.Get("/:companyId/contract", s.ContractHandler.GetActiveContract)
    companies.Get("/:companyId/services/:service/availability", s.ContractHandler.CheckAvailability)
}

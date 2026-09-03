package api

import (
	"job4j_go_share_trip_contracts/gen"
	"job4j_go_share_trip_contracts/internal/domain/contract/repository"
	"job4j_go_share_trip_contracts/internal/domain/contract/service"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	ContractHandler *ContractHandler
}

// ====== Реализация методов gen.ServerInterface ======

func (s *Server) GetContract(c *fiber.Ctx, contractId gen.ContractId) error {
	return s.ContractHandler.GetContract(c, contractId)
}

func (s *Server) CreateContract(c *fiber.Ctx) error {
	return s.ContractHandler.CreateContract(c)
}

func (s *Server) SignContract(c *fiber.Ctx) error {
	return s.ContractHandler.SignContract(c)
}

func (s *Server) GetActiveContract(c *fiber.Ctx, companyId gen.CompanyId) error {
	return s.ContractHandler.GetActiveContract(c, companyId)
}

func (s *Server) ChangeContractStatus(c *fiber.Ctx, contractId gen.ContractId) error {
	return s.ContractHandler.ChangeContractStatus(c, contractId)
}

func (s *Server) UpdateContractServices(c *fiber.Ctx, contractId gen.ContractId) error {
	return s.ContractHandler.UpdateContractServices(c, contractId)
}

func (s *Server) CheckServiceAvailability(c *fiber.Ctx, companyId gen.CompanyId, service gen.Service) error {
	return s.ContractHandler.CheckServiceAvailability(c, companyId, service)
}

func NewServer(pgxpool *pgxpool.Pool) *Server {
	contractRepo := repository.NewRepository(pgxpool)
	contractSvc := service.NewContractService(contractRepo)
	contractHandler := NewContractHandler(contractSvc)

	return &Server{
		ContractHandler: contractHandler,
	}
}
// internal/api/server.go
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

// GetContract Получить договор по идентификатору
// (GET /contracts/{contractId})
func (s *Server) GetContract(c *fiber.Ctx, contractId gen.ContractId) error {
	return s.ContractHandler.GetContract(c, contractId)
}

// CreateContract Создать договор для компании
// (POST /contracts)
func (s *Server) CreateContract(c *fiber.Ctx) error {
	return s.ContractHandler.CreateContract(c)
}

// SignContract Sign a contract
// (POST /contracts/sign_contract)
func (s *Server) SignContract(c *fiber.Ctx) error {
	return s.ContractHandler.SignContract(c)
}

// GetActiveContract Получить активный договор компании
// (GET /companies/{companyId}/contract)
func (s *Server) GetActiveContract(c *fiber.Ctx, companyId gen.CompanyId) error {
	return s.ContractHandler.GetActiveContract(c, companyId)
}

// ChangeContractStatus Изменить статус договора
// (PATCH /contracts/{contractId}/status)
func (s *Server) ChangeContractStatus(c *fiber.Ctx, contractId gen.ContractId) error {
	return s.ContractHandler.ChangeContractStatus(c, contractId)
}

// UpdateContractServices Добавить или обновить список доступных услуг
// (PUT /contracts/{contractId}/services)
func (s *Server) UpdateContractServices(c *fiber.Ctx, contractId gen.ContractId) error {
	return s.ContractHandler.UpdateContractServices(c, contractId)
}

// CheckServiceAvailability Проверить доступность услуги для компании
// (GET /companies/{companyId}/services/{service}/availability)
func (s *Server) CheckServiceAvailability(c *fiber.Ctx, companyId gen.CompanyId, service gen.Service) error {
	return s.ContractHandler.CheckServiceAvailability(c, companyId, service)
}

// ====== Конструктор ======

func NewServer(pgxpool *pgxpool.Pool) *Server {
	// Инициализация контрактов
	contractRepo := repository.NewRepository(pgxpool)
	contractSvc := service.NewContractService(contractRepo)
	contractHandler := NewContractHandler(contractSvc)

	return &Server{
		ContractHandler: contractHandler,
	}
}
package api

import (
	"job4j_go_share_trip_contracts/internal/domain/contract/repository"
	"job4j_go_share_trip_contracts/internal/domain/contract/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct{
    ContractHandler *ContractHandler
}

func NewServer(pgxpool *pgxpool.Pool) *Server {
	// Инициализация контрактов
	contractRepo := repository.NewRepository(pgxpool)

	contractSvc := service.NewContractService(contractRepo)
	contractHandler := NewContractHandler(contractSvc)

	return &Server{
		ContractHandler: contractHandler,
	}
}

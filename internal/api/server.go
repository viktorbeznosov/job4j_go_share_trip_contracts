package api

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct{

}

func NewServer(pgxpool *pgxpool.Pool) *Server {
	return &Server{

	}
}

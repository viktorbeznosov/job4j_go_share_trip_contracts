// internal/api/route.go
package api

import (
	"job4j_go_share_trip_contracts/gen"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router, server *Server) {
	gen.RegisterHandlers(router, server)
}
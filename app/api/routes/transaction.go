package routes

import (
	"github.com/brunomdev/digital-account/app/api/handlers"
	"github.com/gofiber/fiber/v2"
)

func TransactionRoutes(route *fiber.App, handler handlers.TransactionHandler) {
	routes := route.Group("/transactions")
	routes.Post("/", handler.Create)
}

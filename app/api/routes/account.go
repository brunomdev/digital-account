package routes

import (
	"github.com/brunomdev/digital-account/app/api/handlers"
	"github.com/gofiber/fiber/v2"
)

func AccountRoutes(route *fiber.App, handler handlers.AccountHandler) {
	routes := route.Group("/accounts")
	routes.Post("/", handler.Create)
	routes.Get("/:id", handler.Get)
}

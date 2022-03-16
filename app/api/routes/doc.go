package routes

import (
	"github.com/gofiber/fiber/v2"
)

func DocRoutes(route *fiber.App) {
	routes := route.Group("/docs")
	routes.Static("/", "./docs")
}

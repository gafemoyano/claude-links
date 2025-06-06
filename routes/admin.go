package routes

import (
	"be-links/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupAdminRoutes(app *fiber.App, linkHandler *handlers.LinkHandler) {
	admin := app.Group("/admin")

	admin.Post("/create", linkHandler.CreateLink)
	admin.Get("/links", linkHandler.GetAllLinks)
}
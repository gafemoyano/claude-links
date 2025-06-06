package routes

import (
	"be-links/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

func SetupAdminRoutes(app *fiber.App, linkHandler *handlers.LinkHandler) {
	admin := app.Group("/admin")
	
	admin.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": "password", // TODO: Use environment variables
		},
		Realm: "Forbidden",
		Unauthorized: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		},
	}))

	admin.Post("/create", linkHandler.CreateLink)
}
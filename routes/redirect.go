package routes

import (
	"be-links/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRedirectRoutes(app *fiber.App, linkHandler *handlers.LinkHandler) {
	app.Get("/:shortcode", linkHandler.RedirectLink)
	app.Get("/info/:shortcode", linkHandler.GetLinkInfo)
	app.Get("/deeplink/:shortcode", linkHandler.GetDeepLink)
}
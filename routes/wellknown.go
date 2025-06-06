package routes

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func SetupWellKnownRoutes(app *fiber.App) {
	// Dynamic .well-known files
	app.Get("/.well-known/assetlinks.json", func(c *fiber.Ctx) error {
		// Get environment variables for Android app configuration
		packageName := os.Getenv("ANDROID_PACKAGE_NAME")
		appName := os.Getenv("APP_NAME")
		sha256Fingerprint := os.Getenv("ANDROID_SHA256_FINGERPRINT")

		// Set defaults for dev environment
		if packageName == "" {
			packageName = "com.trii.qa"
		}
		if sha256Fingerprint == "" {
			sha256Fingerprint = "SHA256_FINGERPRINT_HERE"
		}

		assetlinks := []map[string]any{
			{
				"relation": []string{"delegate_permission/common.handle_all_urls"},
				"target": map[string]any{
					"namespace":                appName,
					"package_name":             packageName,
					"sha256_cert_fingerprints": []string{sha256Fingerprint},
				},
			},
		}

		c.Set("Content-Type", "application/json")
		return c.JSON(assetlinks)
	})

	app.Get("/.well-known/apple-app-site-association", func(c *fiber.Ctx) error {
		// Get environment variables for iOS app configuration
		iosAppId := os.Getenv("IOS_APP_ID")

		// Set default for dev environment
		if iosAppId == "" {
			iosAppId = "com.trii.qa"
		}

		appleAppSiteAssociation := map[string]any{
			"applinks": map[string]any{
				"apps": []string{},
				"details": []map[string]any{
					{
						"appID": iosAppId,
						"paths": []string{"*"},
					},
				},
			},
		}

		c.Set("Content-Type", "application/json")
		return c.JSON(appleAppSiteAssociation)
	})
}

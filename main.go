package main

import (
	"log"
	"os"

	"be-links/handlers"
	"be-links/routes"
	"be-links/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	_ "be-links/docs" // Import generated docs
)

// @title Dynamic Link Service API
// @version 1.0
// @description This service replaces Firebase Dynamic Links. It supports generating short URLs that redirect users to appropriate destinations based on device platform.
// @host localhost:3000
// @BasePath /
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, proceeding with defaults")
	}

	// Environment variables expected:
	// DATABASE_URL - PostgreSQL connection string
	// IOS_STORE_URL - iOS App Store URL (optional, defaults to https://apps.apple.com/app/id123456789)
	// ANDROID_STORE_URL - Android Play Store URL (optional, defaults to https://play.google.com/store/apps/details?id=com.triico.app&hl=en)
	// PORT - Server port (optional, defaults to 3000)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(logger.New())
	app.Use(cors.New())

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// Serve static files from public directory
	app.Static("/", "./public")

	// Specifically serve .well-known files with correct Content-Type
	app.Static("/.well-known", "./public/.well-known", fiber.Static{
		Compress: false,
		Browse:   false,
	})

	// Initialize database
	db, err := storage.NewDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize handlers
	linkHandler := handlers.NewLinkHandler(db)

	// Swagger documentation
	app.Get("/docs/*", fiberSwagger.WrapHandler)

	// Setup routes
	routes.SetupRedirectRoutes(app, linkHandler)
	routes.SetupAdminRoutes(app, linkHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

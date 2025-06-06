package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"be-links/models"
	"be-links/storage"

	"github.com/gofiber/fiber/v2"
)

type LinkHandler struct {
	db *storage.DB
}

func NewLinkHandler(db *storage.DB) *LinkHandler {
	return &LinkHandler{db: db}
}

func (h *LinkHandler) CreateLink(c *fiber.Ctx) error {
	var req models.CreateLinkRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.DeepLink == "" || req.IOSStore == "" || req.AndroidStore == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "deep_link, ios_store, and android_store are required",
		})
	}

	shortCode, err := generateShortCode()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate short code",
		})
	}

	now := time.Now()
	link := &models.Link{
		ID:           shortCode,
		DeepLink:     req.DeepLink,
		IOSStore:     req.IOSStore,
		AndroidStore: req.AndroidStore,
		Title:        req.Title,
		Description:  req.Description,
		CreatedAt:    now,
		UpdatedAt:    now,
		ClickCount:   0,
	}

	if err := h.db.CreateLink(link); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create link",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(link)
}

func (h *LinkHandler) RedirectLink(c *fiber.Ctx) error {
	shortCode := c.Params("shortcode")
	if shortCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Short code is required",
		})
	}

	link, err := h.db.GetLink(shortCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	if link == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Link not found",
		})
	}

	go h.db.IncrementClickCount(shortCode)

	userAgent := c.Get("User-Agent")
	platform := detectPlatform(userAgent)

	switch platform {
	case "ios":
		return c.Redirect(link.IOSStore, fiber.StatusFound)
	case "android":
		intentURL := buildAndroidIntent(link.DeepLink, link.AndroidStore)
		return c.Redirect(intentURL, fiber.StatusFound)
	default:
		// For unknown platforms, try the deep link first
		// If the app is installed, it will intercept the deep link
		// If not, the browser will hit our server again and we can redirect to a store
		return c.Redirect(link.DeepLink, fiber.StatusFound)
	}
}

func (h *LinkHandler) GetLinkInfo(c *fiber.Ctx) error {
	shortCode := c.Params("shortcode")
	if shortCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Short code is required",
		})
	}

	link, err := h.db.GetLink(shortCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	if link == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Link not found",
		})
	}

	return c.JSON(link)
}

func generateShortCode() (string, error) {
	bytes := make([]byte, 3)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func detectPlatform(userAgent string) string {
	userAgent = strings.ToLower(userAgent)
	
	if strings.Contains(userAgent, "iphone") || strings.Contains(userAgent, "ipad") || strings.Contains(userAgent, "ipod") {
		return "ios"
	}
	
	if strings.Contains(userAgent, "android") {
		return "android"
	}
	
	return "unknown"
}

func buildAndroidIntent(deepLink, fallbackURL string) string {
	return "intent://" + strings.TrimPrefix(deepLink, "myapp://") + 
		"#Intent;scheme=myapp;package=com.example.app;S.browser_fallback_url=" + 
		fallbackURL + ";end"
}
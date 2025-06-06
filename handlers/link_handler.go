package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"os"
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

// CreateLink creates a new short link
// @Summary Create a new short link
// @Description Create a new short link that redirects to the specified deep link
// @Tags links
// @Accept json
// @Produce json
// @Param link body models.CreateLinkRequest true "Link creation request"
// @Success 201 {object} models.Link
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /create [post]
func (h *LinkHandler) CreateLink(c *fiber.Ctx) error {
	var req models.CreateLinkRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.UniversalLink == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "universal_link is required",
		})
	}

	if req.IOSStore == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ios_store is required",
		})
	}

	if req.AndroidStore == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "android_store is required",
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
		ID:            shortCode,
		UniversalLink: req.UniversalLink,
		IOSStore:      req.IOSStore,
		AndroidStore:  req.AndroidStore,
		Title:         req.Title,
		Description:   req.Description,
		CreatedAt:     now,
		UpdatedAt:     now,
		ClickCount:    0,
	}

	if err := h.db.CreateLink(link); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create link",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(link)
}

// RedirectLink redirects to the appropriate platform-specific link using universal/app links
// @Summary Redirect to platform-specific link
// @Description Redirects users using modern universal links (iOS) and app links (Android) based on their User-Agent
// @Tags links
// @Param shortcode path string true "Short code of the link"
// @Success 302 "Redirect to the appropriate link"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /{shortcode} [get]
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
		// For iOS, redirect to universal link first
		// iOS will automatically handle opening the app if installed
		// or falling back to App Store via the universal link infrastructure
		return c.Redirect(link.UniversalLink, fiber.StatusFound)
	case "android":
		// For Android, redirect to universal link first  
		// Android will automatically handle opening the app if installed
		// or falling back to Play Store via the app link infrastructure
		return c.Redirect(link.UniversalLink, fiber.StatusFound)
	default:
		// For desktop/unknown platforms, show a landing page with all options
		return c.JSON(fiber.Map{
			"universal_link": link.UniversalLink,
			"ios_store":      link.IOSStore,
			"android_store":  link.AndroidStore,
			"title":          link.Title,
			"description":    link.Description,
		})
	}
}

// GetLinkInfo retrieves information about a specific link
// @Summary Get link information
// @Description Get detailed information about a specific short link including metadata and click count
// @Tags links
// @Param shortcode path string true "Short code of the link"
// @Success 200 {object} models.Link
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /info/{shortcode} [get]
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

// GetAllLinks retrieves all links from the database
// @Summary Get all links
// @Description Get a list of all short links in the system
// @Tags links
// @Produce json
// @Success 200 {array} models.Link
// @Failure 500 {object} map[string]string
// @Router /admin/links [get]
func (h *LinkHandler) GetAllLinks(c *fiber.Ctx) error {
	links, err := h.db.GetAllLinks()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve links",
		})
	}

	return c.JSON(links)
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


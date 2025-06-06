package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/url"
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
// @Router /admin/create [post]
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

	// Get store URLs from environment variables
	iosStore := os.Getenv("IOS_STORE_URL")
	if iosStore == "" {
		iosStore = "https://apps.apple.com/co/app/trii/id1513826307"
	}

	androidStore := os.Getenv("ANDROID_STORE_URL")
	if androidStore == "" {
		androidStore = "https://play.google.com/store/apps/details?id=com.triico.app&hl=en"
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
		DeepLink:      req.DeepLink,
		IOSStore:      iosStore,
		AndroidStore:  androidStore,
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

	// Add the full short URL to the response
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}
	
	response := fiber.Map{
		"id":             link.ID,
		"short_url":      baseURL + "/" + link.ID,
		"universal_link": link.UniversalLink,
		"deep_link":      link.DeepLink,
		"ios_store":      link.IOSStore,
		"android_store":  link.AndroidStore,
		"title":          link.Title,
		"description":    link.Description,
		"created_at":     link.CreatedAt,
		"updated_at":     link.UpdatedAt,
		"click_count":    link.ClickCount,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
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
		// For desktop/unknown platforms, redirect to trii.co with query parameters
		redirectURL := buildDesktopRedirectURL(link)
		return c.Redirect(redirectURL, fiber.StatusFound)
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

	// Add generated deeplink and short URL to response
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}
	deeplink := generateDeepLink(link.DeepLink, link.UniversalLink)
	response := fiber.Map{
		"id":             link.ID,
		"short_url":      baseURL + "/" + link.ID,
		"universal_link": link.UniversalLink,
		"deep_link":      link.DeepLink,
		"deeplink":       deeplink,
		"ios_store":      link.IOSStore,
		"android_store":  link.AndroidStore,
		"title":          link.Title,
		"description":    link.Description,
		"created_at":     link.CreatedAt,
		"updated_at":     link.UpdatedAt,
		"click_count":    link.ClickCount,
	}

	return c.JSON(response)
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

	// Add generated deeplinks and short URLs to each link in the response
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}
	var enrichedLinks []fiber.Map
	for _, link := range links {
		deeplink := generateDeepLink(link.DeepLink, link.UniversalLink)
		enrichedLink := fiber.Map{
			"id":             link.ID,
			"short_url":      baseURL + "/" + link.ID,
			"universal_link": link.UniversalLink,
			"deep_link":      link.DeepLink,
			"deeplink":       deeplink,
			"ios_store":      link.IOSStore,
			"android_store":  link.AndroidStore,
			"title":          link.Title,
			"description":    link.Description,
			"created_at":     link.CreatedAt,
			"updated_at":     link.UpdatedAt,
			"click_count":    link.ClickCount,
		}
		enrichedLinks = append(enrichedLinks, enrichedLink)
	}

	return c.JSON(enrichedLinks)
}

// GetDeepLink generates a deeplink with triiapp:// scheme for a given short code
// @Summary Get deeplink for short code
// @Description Generate a deeplink with triiapp:// scheme that opens the app directly
// @Tags links
// @Param shortcode path string true "Short code of the link"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /deeplink/{shortcode} [get]
func (h *LinkHandler) GetDeepLink(c *fiber.Ctx) error {
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

	deeplink := generateDeepLink(link.DeepLink, link.UniversalLink)

	return c.JSON(fiber.Map{
		"deeplink":       deeplink,
		"universal_link": link.UniversalLink,
		"deep_link":      link.DeepLink,
		"title":          link.Title,
		"description":    link.Description,
	})
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

func generateDeepLink(deepLink, universalLink string) string {
	// If a custom deep link is provided, use it with triiapp:// scheme
	if deepLink != "" {
		// If it already has a scheme, use as is
		if strings.Contains(deepLink, "://") {
			return deepLink
		}
		// Otherwise, add triiapp:// scheme
		return "triiapp://" + deepLink
	}

	// If no custom deep link, generate one from the universal link
	if universalLink != "" {
		// Extract path from universal link and convert to deeplink
		// Example: https://yourdomain.com/app/product?id=987 -> triiapp://product?id=987
		if strings.Contains(universalLink, "/app/") {
			parts := strings.Split(universalLink, "/app/")
			if len(parts) > 1 {
				return "triiapp://" + parts[1]
			}
		}
	}

	// Fallback to basic deeplink
	return "triiapp://home"
}

func buildDesktopRedirectURL(link *models.Link) string {
	// Get website URL from environment variable
	websiteURL := os.Getenv("WEBSITE_URL")
	if websiteURL == "" {
		websiteURL = "https://trii.co"
	}
	
	// Parse the universal link to extract query parameters
	parsedUniversalLink, err := url.Parse(link.UniversalLink)
	if err != nil {
		return websiteURL
	}
	
	// Create new URL with website base
	redirectURL, err := url.Parse(websiteURL)
	if err != nil {
		return websiteURL
	}
	
	// Copy query parameters from the universal link stored in DB
	redirectURL.RawQuery = parsedUniversalLink.RawQuery
	
	return redirectURL.String()
}

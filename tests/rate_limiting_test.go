package tests

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestRateLimiting(t *testing.T) {
	app := fiber.New()

	// Rate limiting middleware'i ekle
	app.Use(limiter.New(limiter.Config{
		Max:        2, // Maksimum 2 istek
		Expiration: 1 * time.Minute,
	}))

	// Test route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// İlk istek
	req, _ := http.NewRequest(fiber.MethodGet, "/", nil)
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// İkinci istek
	resp, err = app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Üçüncü istek (limit aşılacak)
	resp, err = app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 429, resp.StatusCode) // 429 Too Many Requests
}

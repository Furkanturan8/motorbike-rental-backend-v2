package middleware

import (
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/model"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/pkg/errorx"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/pkg/jwt"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Authorization header kontrolü
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return errorx.WithDetails(errorx.ErrUnauthorized, "Authorization header bulunamadı")
		}

		// Bearer token formatı kontrolü
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return errorx.WithDetails(errorx.ErrInvalidRequest, "Geçersiz Authorization header formatı. 'Bearer <token>' formatında olmalı")
		}

		// Token doğrulama
		claims, err := jwt.Validate(tokenParts[1])
		if err != nil {
			return errorx.WithDetails(errorx.ErrUnauthorized, "Geçersiz veya süresi dolmuş token")
		}

		// Context'e kullanıcı bilgilerini ekle
		c.Locals("userID", claims.UserID)
		c.Locals("role", claims.Role)
		c.Locals("status", claims.Status)
		c.Locals("email", claims.Email)

		return c.Next()
	}
}

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil {
			return errorx.WithDetails(errorx.ErrUnauthorized, "Yetkilendirme bilgisi bulunamadı")
		}

		if role.(model.Role) != model.AdminRole {
			return errorx.WithDetails(errorx.ErrForbidden, "Bu işlem için admin yetkisi gerekli")
		}

		return c.Next()
	}
}

// Belirli rollere sahip kullanıcılar için middleware
func HasRole(roles ...model.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil {
			return errorx.WithDetails(errorx.ErrUnauthorized, "Yetkilendirme bilgisi bulunamadı")
		}

		userRole := role.(model.Role)
		for _, allowedRole := range roles {
			if userRole == allowedRole {
				return c.Next()
			}
		}

		return errorx.WithDetails(errorx.ErrForbidden, "Bu işlem için yeterli yetkiniz yok!")

	}
}

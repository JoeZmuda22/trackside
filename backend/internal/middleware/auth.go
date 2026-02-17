package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
}

func AuthMiddleware(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			}

			tokenStr := parts[1]
			claims := &JWTClaims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			}

			c.Set("userId", claims.Subject)
			c.Set("email", claims.Email)
			return next(c)
		}
	}
}

// GetUserID extracts the user ID from the Echo context (set by AuthMiddleware).
func GetUserID(c echo.Context) string {
	if v, ok := c.Get("userId").(string); ok {
		return v
	}
	return ""
}

// GetUserEmail extracts the user email from the Echo context.
func GetUserEmail(c echo.Context) string {
	if v, ok := c.Get("email").(string); ok {
		return v
	}
	return ""
}

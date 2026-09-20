package middleware

import (
	"avto-crm-api/internal/utils"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type APIResponse struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Data interface{} `json:"data,omitempty"`
	Error interface{} `json:"error,omitempty"`
	Meta interface{} `json:"meta,omitempty"`
}

func AuthMiddleware(secret string, isBlacklist func(ctx context.Context, token string) (bool, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
				Success: true,
				Message: "not authorized",
				Error: utils.ErrUnauthorized,
			})
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
				Success: true,
				Message: "invalid authorized header",
				Error: utils.ErrUnauthorized,
			})
			return
		}

		tokenString := strings.TrimSpace(parts[1])

		token, err := jwt.ParseWithClaims(
			tokenString,
			&Claims{},
			func(t *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			},
		)

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
				Success: true,
				Message: "invalid token",
				Error: utils.ErrUnauthorized,
			})
			return
		}

		// tokenOk, err := isBlacklist(c.Request.Context(), tokenString)

		// if err != nil {
		// 	c.AbortWithStatusJSON(http.StatusInternalServerError, APIResponse{
		// 		Success: true,
		// 		Message: err.Error(),
		// 		Error: err,
		// 	})
		// 	return
		// }

		// if tokenOk {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
		// 		Success: true,
		// 		Message: "not authorized",
		// 		Error: utils.ErrUnauthorized,
		// 	})
		// 	return
		// }

		claims, ok := token.Claims.(*Claims)
		if !ok || claims.UserID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
				Success: true,
				Message: "invalid claims",
				Error: utils.ErrUnauthorized,
			})
			return
		}

		c.Set("userId", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("token", token)

		c.Next()
	}
}

func AdminMiddleware(secret string, isBlacklist func(ctx context.Context, token string) (bool, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
				Success: true,
				Message: "not authorized",
				Error: utils.ErrUnauthorized,
			})
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
				Success: true,
				Message: "invalid authorized header",
				Error: utils.ErrUnauthorized,
			})
			return
		}

		tokenString := strings.TrimSpace(parts[1])

		token, err := jwt.ParseWithClaims(
			tokenString,
			&Claims{},
			func(t *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			},
		)

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
				Success: true,
				Message: "invalid token",
				Error: utils.ErrUnauthorized,
			})
			return
		}

		tokenOk, err := isBlacklist(c.Request.Context(), tokenString)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, APIResponse{
				Success: true,
				Message: err.Error(),
				Error: err,
			})
			return
		}

		if tokenOk {
			c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
				Success: true,
				Message: "not authorized",
				Error: utils.ErrUnauthorized,
			})
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok || claims.UserID == "" || claims.Role != "admin" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
				Success: true,
				Message: "invalid claims",
				Error: utils.ErrUnauthorized,
			})
			return
		}

		c.Set("userId", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("token", token)

		c.Next()
	}
}

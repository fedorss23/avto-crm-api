package middleware

import (
	"avto-crm-api/internal/modules/company"
	"avto-crm-api/internal/utils"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type ManagerClaims struct {
	ManagerID string `json:"manager_id"`
	Role      string `json:"role"`
	CompanyID string `json:"company_id"`
	jwt.RegisteredClaims
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func checkTokenAndUserClaims(c *gin.Context, secret string) (*UserClaims, bool) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, false
	}

	parts := strings.SplitN(authHeader, " ", 2)

	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, false
	}

	tokenString := strings.TrimSpace(parts[1])

	token, err := jwt.ParseWithClaims(
		tokenString,
		&UserClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unsupported signing method")
			}
			return []byte(secret), nil
		},
	)

	if err != nil || !token.Valid {
		return nil, false
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || claims.UserID == "" {
		return nil, false
	}

	return claims, true
}

func checkTokenAndManagerClaims(c *gin.Context, secret string) (*ManagerClaims, bool) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, false
	}

	parts := strings.SplitN(authHeader, " ", 2)

	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, false
	}

	tokenString := strings.TrimSpace(parts[1])

	token, err := jwt.ParseWithClaims(
		tokenString,
		&ManagerClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unsupported signing method")
			}
			return []byte(secret), nil
		},
	)

	if err != nil || !token.Valid {
		return nil, false
	}

	claims, ok := token.Claims.(*ManagerClaims)
	if !ok || claims.ManagerID == "" {
		return nil, false
	}

	return claims, true
}

func UserAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		claims, ok := checkTokenAndUserClaims(c, secret)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Message: "invalid token",
				Error:   utils.ErrUnauthorized,
			})
			return
		}

		c.Set("userId", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("type", "owner")
		c.Next()
	}
}

func UserAuthMiddlewareWithCompany(secret string, findCompany func(ownerId, companyId uuid.UUID) (*company.Company, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		claims, ok := checkTokenAndUserClaims(c, secret)

		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Message: "invalid token",
				Error:   utils.ErrUnauthorized,
			})
			return
		}

		companyId, exists := c.GetQuery("companyId")
		if !exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, APIResponse{
				Success: false,
				Message: "validation error",
				Data: utils.ValidationErrorResponseData{
					Err:  map[string]string{"companyId": "required"},
					Code: utils.ValidationErrorCode,
				},
			})
			return
		}

		log.Println(claims.UserID)
		log.Println(companyId)

		oid, err := uuid.Parse(claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, APIResponse{
				Success: false,
				Message: "userId error",
				Data: utils.ErrorResponseData{
					Err:  errors.New("invalid user id").Error(),
					Code: utils.ValidationErrorCode,
				},
			})
		}

		cid, err := uuid.Parse(companyId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, APIResponse{
				Success: false,
				Message: "companyId error",
				Data: utils.ErrorResponseData{
					Err:  errors.New("invalid company id").Error(),
					Code: utils.ValidationErrorCode,
				},
			})
		}

		if _, err := findCompany(oid, cid); err != nil {
			c.AbortWithStatusJSON(utils.ErrorToHTTPStatus(err), APIResponse{
				Success: false,
				Message: "company not found",
				Data: utils.ErrorResponseData{
					Err:  utils.ErrRecordNotFound.Error(),
					Code: utils.CodeByError(err),
				},
			})
			return
		}

		c.Set("userId", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("type", "owner")
		c.Set("companyId", companyId)
		c.Next()
	}
}

func ManagerAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		claims, ok := checkTokenAndManagerClaims(c, secret)

		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Message: "invalid token",
				Error:   utils.ErrUnauthorized,
			})
			return
		}

		c.Set("managerId", claims.ManagerID)
		c.Set("role", claims.Role)
		c.Set("companyId", claims.CompanyID)
		c.Set("type", "manager")
		c.Next()
	}
}

func StandartMiddlewareWithCompany(secret string, findCompany func(ownerId, companyId uuid.UUID) (*company.Company, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		claims, ok := checkTokenAndUserClaims(c, secret)
		if !ok {
			mclaims, ok := checkTokenAndManagerClaims(c, secret)
			log.Println(mclaims.ManagerID, ok)

			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
					Success: false,
					Message: "invalid token",
					Error:   utils.ErrUnauthorized,
				})
				return
			}
			c.Set("userId", mclaims.ManagerID)
			c.Set("role", mclaims.Role)
			c.Set("companyId", mclaims.CompanyID)
			c.Set("type", "manager")
			c.Next()
			return
		}

		companyId, exists := c.GetQuery("companyId")
		if !exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, APIResponse{
				Success: false,
				Message: "validation error",
				Data: utils.ValidationErrorResponseData{
					Err:  map[string]string{"companyId": "required"},
					Code: utils.ValidationErrorCode,
				},
			})
			return
		}

		log.Println(claims.UserID)
		log.Println(companyId)

		oid, err := uuid.Parse(claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, APIResponse{
				Success: false,
				Message: "userId error",
				Data: utils.ErrorResponseData{
					Err:  errors.New("invalid user id").Error(),
					Code: utils.ValidationErrorCode,
				},
			})
		}

		cid, err := uuid.Parse(companyId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, APIResponse{
				Success: false,
				Message: "companyId error",
				Data: utils.ErrorResponseData{
					Err:  errors.New("invalid company id").Error(),
					Code: utils.ValidationErrorCode,
				},
			})
		}

		if _, err := findCompany(oid, cid); err != nil {
			c.AbortWithStatusJSON(utils.ErrorToHTTPStatus(err), APIResponse{
				Success: false,
				Message: "company not found",
				Data: utils.ErrorResponseData{
					Err:  utils.ErrRecordNotFound.Error(),
					Code: utils.CodeByError(err),
				},
			})
			return
		}

		c.Set("userId", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("type", "owner")
		c.Set("companyId", companyId)
		c.Next()
	}
}

func StandartMiddleware(secret string, findCompany func(ownerId, companyId uuid.UUID) (*company.Company, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		claims, ok := checkTokenAndUserClaims(c, secret)
		if !ok {
			mclaims, ok := checkTokenAndManagerClaims(c, secret)
			log.Println(mclaims.ManagerID, ok)

			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
					Success: false,
					Message: "invalid token",
					Error:   utils.ErrUnauthorized,
				})
				return
			}
			c.Set("userId", mclaims.ManagerID)
			c.Set("role", mclaims.Role)
			c.Set("companyId", mclaims.CompanyID)
			c.Set("type", "manager")
			c.Next()
			return
		}

		companyId, exists := c.GetQuery("companyId")

		log.Println(claims.UserID)
		log.Println(companyId)

		oid, err := uuid.Parse(claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, APIResponse{
				Success: false,
				Message: "userId error",
				Data: utils.ErrorResponseData{
					Err:  errors.New("invalid user id").Error(),
					Code: utils.ValidationErrorCode,
				},
			})
		}

		if exists && companyId != "" {
			cid, err := uuid.Parse(companyId)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, APIResponse{
					Success: false,
					Message: "companyId error",
					Data: utils.ErrorResponseData{
						Err:  errors.New("invalid company id").Error(),
						Code: utils.ValidationErrorCode,
					},
				})
			}

			if _, err := findCompany(oid, cid); err != nil {
				c.AbortWithStatusJSON(utils.ErrorToHTTPStatus(err), APIResponse{
					Success: false,
					Message: "company not found",
					Data: utils.ErrorResponseData{
						Err:  utils.ErrRecordNotFound.Error(),
						Code: utils.CodeByError(err),
					},
				})
				return
			}
			c.Set("companyId", companyId)
		}

		c.Set("userId", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("type", "owner")
		c.Next()
	}
}

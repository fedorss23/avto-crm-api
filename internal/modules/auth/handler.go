package auth

import (
	"avto-crm-api/internal/blacklist"
	"avto-crm-api/internal/utils"
	"avto-crm-api/pkg/cookie"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service          *AuthService
	cookieConfig     *cookie.CookieConfig
	blacklistService *blacklist.BlacklistService
}

func NewAuthHandler(service *AuthService, cfg *cookie.CookieConfig, blacklistService *blacklist.BlacklistService) *AuthHandler {
	return &AuthHandler{
		service:          service,
		cookieConfig:     cfg,
		blacklistService: blacklistService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errs := utils.ParseValidationErrors(err)
		utils.ValidationErrorResponse(c, errs)
		return
	}

	authresp, err := h.service.Register(&req)
	if err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), err.Error(), err, utils.CodeByError(err))
		return
	}

	h.cookieConfig.SetRefreshToken(c, authresp.RefreshToken)

	ans := &RegisterResponse{
		User:        authresp.User,
		AccessToken: authresp.AccessToken,
	}

	utils.SuccessResponse(c, http.StatusCreated, "Registration successful, please log in", ans)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errs := utils.ParseValidationErrors(err)
		utils.ValidationErrorResponse(c, errs)
		return
	}

	authresp, err := h.service.Login(&req, c.ClientIP())

	if err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), err.Error(), err, utils.CodeByError(err))
		return
	}

	h.cookieConfig.SetRefreshToken(c, authresp.RefreshToken)

	ans := &RegisterResponse{
		User:        authresp.User,
		AccessToken: authresp.AccessToken,
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", ans)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, refreshTtl, err := h.getTtlForRefreshToken(c)
	if err != nil {
		return
	}

	tokens, err := h.service.RefreshToken(refreshToken)

	if err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), err.Error(), err, utils.CodeByError(err))
		return
	}

	h.cookieConfig.SetRefreshToken(c, tokens.RefreshToken)

	accessToken, accessTtl, err := h.getTtlForAccessToken(c)
	if err != nil {
		return
	}

	if err := h.blacklistService.Add(c.Request.Context(), refreshToken, refreshTtl); err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), err.Error(), err, utils.CodeByError(err))
		return
	}

	if err := h.blacklistService.Add(c.Request.Context(), accessToken, accessTtl); err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), err.Error(), err, utils.CodeByError(err))
		return
	}

	ans := &RegisterResponse{
		User:        tokens.User,
		AccessToken: tokens.AccessToken,
	}

	utils.SuccessResponse(c, http.StatusOK, "Token successfully refreshed", ans)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// refreshToken, refreshTtl, err := h.getTtlForRefreshToken(c)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// accessToken, accessTtl, err := h.getTtlForAccessToken(c)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// if err := h.blacklistService.Add(c.Request.Context(), refreshToken, refreshTtl); err != nil {
	// 	fmt.Println(err)
	// 	utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), err.Error(), err, utils.CodeByError(err))
	// 	return
	// }

	// if err := h.blacklistService.Add(c.Request.Context(), accessToken, accessTtl); err != nil {
	// 	fmt.Println(err)
	// 	utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), err.Error(), err, utils.CodeByError(err))
	// 	return
	// }

	h.cookieConfig.ClearAuthCookies(c)

	utils.SuccessResponseWithoutBody(c, http.StatusOK, "Logout successful")
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var userId string
	if err := utils.GetOwnerId(c, &userId); err != nil {
		return
	}

	var req ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errs := utils.ParseValidationErrors(err)
		utils.ValidationErrorResponse(c, errs)
		return
	}

	refreshToken, refreshTtl, err := h.getTtlForRefreshToken(c)
	if err != nil {
		return
	}

	accessToken, accessTtl, err := h.getTtlForAccessToken(c)
	if err != nil {
		return
	}

	data, err := h.service.ChangePassword(userId, refreshToken, &req)

	if err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), err.Error(), err, utils.CodeByError(err))
		return
	}

	if err := h.blacklistService.Add(c.Request.Context(), refreshToken, refreshTtl); err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), err.Error(), err, utils.CodeByError(err))
		return
	}

	if err := h.blacklistService.Add(c.Request.Context(), accessToken, accessTtl); err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), err.Error(), err, utils.CodeByError(err))
		return
	}

	h.cookieConfig.ClearAuthCookies(c)

	utils.SuccessResponse(c, http.StatusOK, "Password successfully changed", data)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	var userId string
	if err := utils.GetOwnerId(c, &userId); err != nil {
		return
	}

	user, err := h.service.GetProfile(userId)

	if err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), err.Error(), err, utils.CodeByError(err))
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile found", user)
}

func (h *AuthHandler) CheckAuth(c *gin.Context) {
	var userId string
	if err := utils.GetOwnerId(c, &userId); err != nil {
		return
	}

	user, err := h.service.GetProfile(userId)

	if err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), err.Error(), err, utils.CodeByError(err))
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Authenticated", gin.H{
		"authenticated": true,
		"user":          user,
	})
}

func (h *AuthHandler) getTtlForRefreshToken(c *gin.Context) (string, time.Duration, error) {
	refreshToken, err := h.cookieConfig.GetRefreshToken(c)

	if refreshToken == "" || err != nil {
		utils.ErrorResponse(c,
			utils.ErrorToHTTPStatus(utils.ErrUnauthorized),
			utils.ErrUnauthorized.Error(),
			utils.ErrUnauthorized,
			utils.CodeByError(utils.ErrUnauthorized),
		)
		return "", 0, utils.ErrUnauthorized
	}

	claims, err := h.service.jwtMaker.ValidateRefreshToken(refreshToken)

	if err != nil {
		utils.ErrorResponse(c,
			utils.ErrorToHTTPStatus(utils.ErrUnauthorized),
			utils.ErrUnauthorized.Error(),
			utils.ErrUnauthorized,
			utils.CodeByError(utils.ErrUnauthorized),
		)
		return "", 0, utils.ErrUnauthorized
	}

	return refreshToken, time.Until(claims.ExpiresAt.Time), nil
}

func (h *AuthHandler) getTtlForAccessToken(c *gin.Context) (string, time.Duration, error) {
	accessToken := c.GetString("accessToken")

	if accessToken == "" {
		utils.ErrorResponse(c,
			utils.ErrorToHTTPStatus(utils.ErrUnauthorized),
			utils.ErrUnauthorized.Error(),
			utils.ErrUnauthorized,
			utils.CodeByError(utils.ErrUnauthorized),
		)
		return "", 0, utils.ErrUnauthorized
	}

	claims, err := h.service.jwtMaker.ValidateAccessToken(accessToken)

	if err != nil {
		utils.ErrorResponse(c,
			utils.ErrorToHTTPStatus(utils.ErrUnauthorized),
			utils.ErrUnauthorized.Error(),
			utils.ErrUnauthorized,
			utils.CodeByError(utils.ErrUnauthorized),
		)
		return "", 0, utils.ErrUnauthorized
	}

	return accessToken, time.Until(claims.ExpiresAt.Time), nil
}

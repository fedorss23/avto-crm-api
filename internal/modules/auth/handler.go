package auth

import (
	"avto-crm-api/internal/utils"
	"avto-crm-api/pkg/cookie"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service      *AuthService
	cookieConfig *cookie.CookieConfig
}

func NewAuthHandler(service *AuthService, cfg *cookie.CookieConfig) *AuthHandler {
	return &AuthHandler{
		service:      service,
		cookieConfig: cfg,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errs := utils.ParseValidationErrors(err)
		utils.ValidationErrorResponse(c, errs, utils.ValidationErrorCode)
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
		utils.ValidationErrorResponse(c, errs, utils.ValidationErrorCode)
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
	refreshToken, err := h.cookieConfig.GetRefreshToken(c)

	if err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), err.Error(), err, utils.CodeByError(err))
		return
	}

	tokens, err := h.service.RefreshToken(refreshToken)

	if err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), err.Error(), err, utils.CodeByError(err))
		return
	}

	h.cookieConfig.SetRefreshToken(c, refreshToken)

	ans := &RegisterResponse{
		User:        tokens.User,
		AccessToken: tokens.AccessToken,
	}

	utils.SuccessResponse(c, http.StatusOK, "Token successfully refreshed", ans)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var userId string
	if err := utils.GetOwnerId(c, &userId); err != nil {
		return
	}

	// логика добавления токена в blacklist
	// refreshToken, _ := h.cookieConfig.GetRefreshToken(c)
	// ctx := c.Request.Context()
	// refreshToken из cookie и вызов h.service.Logout(ctx, userID token)
	// далее обработка ошибки если есть

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
		utils.ValidationErrorResponse(c, errs, utils.ValidationErrorCode)
		return
	}

	if err := h.service.ChangePassword(userId, &req); err != nil {
		utils.ErrorResponse(c, utils.ErrorToHTTPStatus(err), err.Error(), err, utils.CodeByError(err))
		return
	}

	h.cookieConfig.ClearAuthCookies(c)

	utils.SuccessResponseWithoutBody(c, http.StatusOK, "Password successfully changed. Please log in again")
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
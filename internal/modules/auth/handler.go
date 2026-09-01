package auth

import (
	"avto-crm-api/internal/utils"
	"avto-crm-api/pkg/cookie"
	"errors"
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
		utils.ValidationErrorResponse(c, errs)
		return
	}

	authresp, err := h.service.Register(&req)
	if err != nil {
		if err == ErrEmailAlreadyExists {
			utils.ErrorResponse(c, http.StatusConflict, "пользователь с таким email уже существует", err)
			return
		}

		utils.ErrorResponse(c, http.StatusInternalServerError, "ошибка сервера", err)
		return
	}

	h.cookieConfig.SetRefreshToken(c, authresp.RefreshToken)

	ans := &RegisterResponse{
		User:        authresp.User,
		AccessToken: authresp.AccessToken,
	}

	utils.SuccessResponse(c, http.StatusCreated, "Успешная регистрация, пожалуйста, войдите", ans)
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
		if err == ErrUserLocked {
			utils.ErrorResponse(c, http.StatusForbidden, "пользователь заблокирован", err)
			return
		}

		utils.ErrorResponse(c, http.StatusBadRequest, "неправильный запрос", err)
		return
	}

	h.cookieConfig.SetRefreshToken(c, authresp.RefreshToken)

	ans := &RegisterResponse{
		User: authresp.User,
		AccessToken: authresp.AccessToken,
	}

	utils.SuccessResponse(c, http.StatusOK, "успешный вход", ans)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := h.cookieConfig.GetRefreshToken(c)

	if err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "пользователь не авторизован", err)
		return
	}

	tokens, err := h.service.RefreshToken(refreshToken)

	if err != nil {
		utils.ErrorResponse(c, ErrorToHTTPStatus(err), "ошибка при обновлении токена", err)
		return
	}

	h.cookieConfig.SetRefreshToken(c, refreshToken)

	ans := &RegisterResponse{
		User: tokens.User,
		AccessToken: tokens.AccessToken,
	}

	utils.SuccessResponse(c, http.StatusOK, "токен успешно обновлен", ans)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID := c.GetString("userId")
	
	if userID == "" {
		utils.ErrorResponse(c, http.StatusForbidden, "пользователь не авторизован", errors.New("wrong get userId as query param"))
		return
	}

	// логика добавления токена в blacklist
	// refreshToken, _ := h.cookieConfig.GetRefreshToken(c)
	// ctx := c.Request.Context()
	// refreshToken из cookie и вызов h.service.Logout(ctx, userID token)
	// далее обработка ошибки если есть

	h.cookieConfig.ClearAuthCookies(c)

	utils.SuccessResponseWithoutBody(c, http.StatusOK, "Успешный выход")
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userId, exists := c.Get("userId")

	if !exists {
		utils.ErrorResponse(c, http.StatusForbidden, "пользователь не авторизован", errors.New("wrong get userId as query param"))
		return
	}

	var req ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errs := utils.ParseValidationErrors(err)
		utils.ValidationErrorResponse(c, errs)
		return
	}

	if err := h.service.ChangePassword(userId.(string), &req); err != nil {
		if err == ErrOldPasswordIncorrect {
			utils.ErrorResponse(c, http.StatusBadRequest, "неправильный пароль", err)
			return
		}

		utils.ErrorResponse(c, http.StatusInternalServerError, "ошибка сервера", err)
		return
	}

	h.cookieConfig.ClearAuthCookies(c)

	utils.SuccessResponseWithoutBody(c, http.StatusOK, "Пароль успешно изменен. Пожалуйста, авторизуйтесь заново")
}


func (h *AuthHandler) GetProfile(c *gin.Context) {
	userId := c.GetString("userId")

	if userId == "" {
		utils.ErrorResponse(c, http.StatusForbidden, "пользователь не авторизован", errors.New("wrong get userId as query param"))
		return
	}

	user, err := h.service.GetProfile(userId)

	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "пользователь не найден", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "профиль найден", user)
}

func (h *AuthHandler) CheckAuth(c *gin.Context) {
	userId, exists := c.Get("userId")

	if !exists {
		utils.ErrorResponse(c, http.StatusForbidden, "пользователь не авторизован", errors.New("wrong get userId as query param"))
		return
	}

	user, err := h.service.GetProfile(userId.(string))

	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "пользователь не авторизован", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "авторизован", gin.H{
		"authenticated": true,
		"user": user,
	})
}
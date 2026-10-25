package authmanagers

import (
	"avto-crm-api/internal/utils"
	"avto-crm-api/pkg/cookie"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthManagersHandler struct {
	authMService *AuthManagersService
	cookieConfig *cookie.CookieConfig
}

func NewAuthManagersHandler(authManagersService *AuthManagersService, cookieConfig *cookie.CookieConfig) *AuthManagersHandler {
	return &AuthManagersHandler{
		authMService: authManagersService,
		cookieConfig: cookieConfig,
	}
}


func (h *AuthManagersHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errs := utils.ParseValidationErrors(err)
		utils.ValidationErrorResponse(c, errs)
		return
	}

	resp, err := h.authMService.Login(&req)

	if err != nil {
		utils.CallErrrorResponse(c, err)
		return
	}

	h.cookieConfig.SetRefreshToken(c, resp.RefreshToken)

	data := &AuthResponseWithoutRefreshToken{
		AccessToken: resp.AccessToken,
		Manager: resp.Manager,
	}

	utils.SuccessResponse(c, http.StatusOK, "sucessfully login", data)
}

func (h *AuthManagersHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errs := utils.ParseValidationErrors(err)
		utils.ValidationErrorResponse(c, errs)
		return
	}

	var ownerId string
	if err := utils.GetOwnerId(c, &ownerId); err != nil {
		return
	}

	companyId := c.GetString("companyId")

	manager, err := h.authMService.Register(&req, ownerId, companyId)

	if err != nil {
		utils.CallErrrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "successfully register", &RegisterResponse{
		Manager: *manager,
	})
}

func (h *AuthManagersHandler) RefreshTokens(c *gin.Context) {
	token, err := h.cookieConfig.GetRefreshToken(c)

	if err != nil {
		utils.CallErrrorResponse(c, err)
		return
	}

	tokens, manager, err := h.authMService.RefreshToken(token)

	if err != nil {
		utils.CallErrrorResponse(c, err)
		return
	}

	h.cookieConfig.SetRefreshToken(c, tokens.RefreshToken)

	utils.SuccessResponse(c, http.StatusOK, "successfully refresh", &AuthResponseWithoutRefreshToken{
		AccessToken: tokens.AccessToken,
		Manager: *manager,
	})

}

func (h *AuthManagersHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errs := utils.ParseValidationErrors(err)
		utils.ValidationErrorResponse(c, errs)
		return
	}

	var ownerId string
	if err := utils.GetOwnerId(c, &ownerId); err != nil {
		return
	}

	var managerId string
	if err := utils.GetParam(c, "managerId", &managerId); err != nil {
		return
	}

	var companyId string
	if err := utils.GetStringRequiredQuery(c, "companyId", &companyId); err != nil {
		return
	}

	manager, err := h.authMService.ChangePassword(req.OldPassword, req.NewPassword, managerId, ownerId, companyId)

	if err != nil {
		utils.CallErrrorResponse(c, err)
	}

	utils.SuccessResponse(c, http.StatusOK, "password successfully changed", &RegisterResponse{
		Manager: *manager,
	})
}

func (h *AuthManagersHandler) Logout(c *gin.Context) {
	h.cookieConfig.ClearAuthCookies(c)
	utils.SuccessResponseWithoutBody(c, http.StatusOK, "successfully logout")
}
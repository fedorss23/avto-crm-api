package authmanagers

import "github.com/google/uuid"


type RegisterRequest struct {
	Login string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	Name string `json:"name" binding:"required"`
	MiddleName string `json:"middleName" binding:"required"`
	LastName *string `json:"lastName,omitempty"`
	Phone *string `json:"phone,omitempty"`
}

type LoginRequest struct {
	Login string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=8,max=72"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type AuthResponse struct {
	AccessToken string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Manager ManagerResponse `json:"manager"`
}

type AuthResponseWithoutRefreshToken struct {
	AccessToken string `json:"accessToken"`
	Manager ManagerResponse `json:"manager"`
}

type ManagerResponse struct {
	ID uuid.UUID `json:"id"`
	Login string `json:"login"`
	Name string `json:"name"`
	MiddleName string `json:"middleName"`
	LastName *string `json:"lastName"`
	CreatedAt string `json:"createdAt"`
}

type TokensResponse struct {
	AccessToken string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

type RegisterResponse struct {
	Manager ManagerResponse `json:"manager"`
}

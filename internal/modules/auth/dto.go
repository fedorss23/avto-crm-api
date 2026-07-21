package auth

import "github.com/google/uuid"


type RegisterRequest struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	Name string `json:"name" binding:"required"`
	LastName string `json:"lastName" binding:"required"`
	Phone string `json:"phone,omitempty"`
}

type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
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
	TokenType string `json:"tokenType"`
	ExpiresIn int64 `json:"expiresIn"`
	User UserResponse `json:"user"`
}

type UserResponse struct {
	ID uuid.UUID `json:"id"`
	Email string `json:"email"`
	Name string `json:"name"`
	LastName string `json:"lastName"`
	Role string `json:"role"`
	CreatedAt string `json:"createdAt"`
}

type TokensResponse struct {
	AccessToken string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn int64 `json:"expiresIn"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

type LoginResult struct {
	User UserResponse `json:"user"`
	AccessToken string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshResult struct {
	User UserResponse `json:"user"`
	AccessToken string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RegisterResponse struct {
	User UserResponse `json:"user"`
	AccessToken string `json:"accessToken"`
}
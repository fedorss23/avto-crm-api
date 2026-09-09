package auth

import (
	"avto-crm-api/internal/modules/user"
	"avto-crm-api/internal/utils"
	"avto-crm-api/pkg/jwt"
	"time"
)

type Config struct {
	AccessTokenDuration time.Duration
	RefreshTokenDuration time.Duration
	MaxLoginAttempts int
	LockDuration time.Duration
}

type AuthService struct {
	userRepo *user.UserRepository
	jwtMaker *jwt.JWTMaker
	config *Config
}

func NewAuthService(userRepo *user.UserRepository, jwtMaker *jwt.JWTMaker, cfg *Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwtMaker: jwtMaker,
		config: cfg,
	}
}

func (s *AuthService) Register(req *RegisterRequest) (*AuthResponse, error) {
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, utils.ErrEmailAlredyExists
	}

	if !utils.ValidatePassword(req.Password) {
		return nil, utils.ErrPasswordIncorrectRegister
	}

	hash, err := utils.Hash(req.Password)

	if err != nil {
		return nil, err
	}

	user := &user.User{
		Email: req.Email,
		Password: string(hash),
		Name: req.Name,
		LastName: req.LastName,
		Role: "user",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	tokens, err := s.generateTokens(user)

	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken: tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType: "Bearer",
		ExpiresIn: int64(s.config.AccessTokenDuration.Seconds()),
		User: *s.toUserResponse(user),
	}, nil
}

func (s *AuthService) Login(req *LoginRequest, ip string) (*LoginResult, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, utils.ErrInvalidCredentials
	}

	// if err := s.isLocked(user); err != nil {
	// 	return nil, ErrUserLocked
	// }

	match, err := utils.Verify(req.Password, user.Password) 

	if err != nil || !match {
		return  nil, utils.ErrPasswordIncorrectLogin
	}

	// s.resetLoginAttempts(user.ID)

	// обновление последнего входа
	// now := time.Now()
	// user.LastLogin := &now

	tokenResponse, err := s.generateTokens(user)

	if err != nil {
		return nil, err
	}

	return &LoginResult{
		User: *s.toUserResponse(user),
		AccessToken: tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*RefreshResult, error) {
	// добавление в blacklist в будущем, интеграция с redis

	claims, err := s.jwtMaker.ValidateRefreshToken(refreshToken)

	if err != nil {
		return nil, utils.ErrInvalidRefreshToken
	}

	user, err := s.userRepo.FindById(claims.UserID)

	if err != nil || user == nil {
		return nil, err
	}

	tokenResponse, err := s.generateTokens(user)

	if err != nil {
		return nil, err
	}

	return &RefreshResult{
		User: *s.toUserResponse(user),
		AccessToken: tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
	}, nil
}

// func (s *AuthService) Logout(userID string) error {
// 	// логика добавления refreshToken в blacklist
// 	// передаем refreshToken в аргументе (берем из контекста gin.Context)

// }

func (s *AuthService) ChangePassword(userID string, req *ChangePasswordRequest) error {
	user, err := s.userRepo.FindById(userID)
	if err != nil || user == nil {
		return err
	}

	match, err := utils.Verify(req.OldPassword, user.Password)

	if err != nil || !match {
		return utils.ErrPasswordIncorrectRegister
	}

	hashPassword, err := utils.Hash(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashPassword

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	//добавление токена в blacklist

	return nil
}

func (s *AuthService) GetProfile(userID string) (*UserResponse, error) {
	user, err := s.userRepo.FindById(userID)

	if err != nil || user == nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

func (s *AuthService) generateTokens(user *user.User) (*TokensResponse, error) {
	accessToken, err := s.jwtMaker.CreateAccessToken(user.ID.String(), user.Email)

	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtMaker.CreateRefreshToken(user.ID.String())

	if err != nil {
		return nil, err
	}

	return &TokensResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		ExpiresIn: int64(s.config.AccessTokenDuration.Seconds()),
	}, nil
}


func (s *AuthService) toUserResponse(user *user.User) *UserResponse {
	return &UserResponse{
		ID: user.ID,
		Email: user.Email,
		Name: user.Name,
		LastName: user.LastName,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}
}

// func (s *AuthService) isLocked(user *user.User) error {
// 	if !user.IsActive {
// 		return ErrUserInactive
// 	}

// 	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
// 		return ErrUserLocked
// 	} 

// 	return nil
// }

// func (s *AuthService) handleFailedLogin(user *user.User) error {
// 	user.LoginAttempts++

// 	if user.LoginAttempts >= s.config.MaxLoginAttempts {
// 		lockedUntil := time.Now().Add(s.config.LockDuration)
// 		user.LockedUntil = &lockedUntil
// 	}

// 	return s.userRepo.Update(user)
// }

// func (s *AuthService) resetLoginAttempts(userID uuid.UUID) error {
// 	return s.userRepo.ResetLoginAttempts(userID)
// }
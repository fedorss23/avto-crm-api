package auth

import (
	"avto-crm-api/internal/modules/user"
	"avto-crm-api/internal/utils"
	"avto-crm-api/pkg/jwt"
	"time"

	"gorm.io/gorm"
	"errors"
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
	db *gorm.DB
}

func NewAuthService(
	userRepo *user.UserRepository, 
	jwtMaker *jwt.JWTMaker, 
	cfg *Config, 
	db *gorm.DB,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwtMaker: jwtMaker,
		config: cfg,
		db: db,
	}
}

func (s *AuthService) Register(req *RegisterRequest) (*AuthResponse, error) {
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err != nil && !errors.Is(err, utils.ErrRecordNotFound){
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

	if user.IsLocked {
		return nil, utils.ErrUserLocked
	}

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

func (s *AuthService) ChangePassword(userID, refreshToken string, req *ChangePasswordRequest) (*AuthResponse, error) {
	var tokens *RefreshResult

	err := s.db.Transaction(func(tx *gorm.DB) error { 
		user, err := s.userRepo.FindByIdWithTx(tx, userID)
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

		if err := s.userRepo.UpdateWithTx(tx, user); err != nil {
			return err
		}

		tokens, err = s.RefreshToken(refreshToken)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User: tokens.User,
		RefreshToken: tokens.RefreshToken,
		AccessToken: tokens.AccessToken,
	}, nil
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

func (s AuthService) LockUser(userId string) error {
	user, err := s.userRepo.FindById(userId)

	if err != nil {
		return err
	}

	user.IsLocked = true

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	return nil
}

// func (s *AuthService) resetLoginAttempts(userID uuid.UUID) error {
// 	return s.userRepo.ResetLoginAttempts(userID)
// }
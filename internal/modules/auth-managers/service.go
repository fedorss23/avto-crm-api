package authmanagers

import (
	"avto-crm-api/internal/modules/company"
	"avto-crm-api/internal/modules/manager"
	"avto-crm-api/internal/utils"
	"avto-crm-api/pkg/jwt"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Config struct {
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	MaxLoginAttempts     int
	LockDuration         time.Duration
}

type AuthManagersService struct {
	managerRepo *manager.ManagerRepository
	companyRepo *company.CompanyRepository
	jwtMaker    *jwt.JWTMaker
	config      *Config
	db          *gorm.DB
}

func NewAuthService(
	managerRepo *manager.ManagerRepository,
	companyRepo *company.CompanyRepository,
	jwtMaker *jwt.JWTMaker,
	cfg *Config,
	db *gorm.DB,
) *AuthManagersService {
	return &AuthManagersService{
		managerRepo: managerRepo,
		companyRepo: companyRepo,
		jwtMaker:    jwtMaker,
		config:      cfg,
		db:          db,
	}
}

func (s *AuthManagersService) Register(req *RegisterRequest, ownerId string, companyId string) (*ManagerResponse, error) {
	var resp *ManagerResponse
	err := s.db.Transaction(func(tx *gorm.DB) error {
		oldManager, err := s.managerRepo.FindByLogin(req.Login)

		if err != nil && !errors.Is(err, utils.ErrRecordNotFound) {
			return err
		}

		if oldManager != nil {
			return utils.ErrLoginAlreadyExists
		}

		hash, err := utils.Hash(req.Password)

		if err != nil {
			return err
		}

		companyIdParse, err := uuid.Parse(companyId)

		if err != nil {
			return err
		}

		newManager := &manager.Manager{
			Login:      req.Login,
			Password:   hash,
			Name:       req.Name,
			MiddleName: req.MiddleName,
			CompanyID:  companyIdParse,
		}

		if req.Phone != nil {
			newManager.Phone = req.Phone
		}

		if req.LastName != nil {
			newManager.LastName = req.LastName
		}

		if err := tx.Create(newManager).Error; err != nil {
			return err
		}

		resp = s.ToManagerResponse(newManager)
		return nil
	})

	return resp, err
}

func (s *AuthManagersService) RefreshToken(refreshToken string) (*TokensResponse, *ManagerResponse, error) {
	claims, err := s.jwtMaker.ValidateRefreshTokenManager(refreshToken)

	if err != nil {
		return nil, nil, utils.ErrInvalidRefreshToken
	}

	manager, err := s.managerRepo.FindById(claims.ManagerID)

	if err != nil {
		return nil, nil, err
	}

	tokens, err := s.generateTokens(manager.ID.String(), manager.Role, manager.CompanyID.String())

	return tokens, s.ToManagerResponse(manager), err
}

func (s *AuthManagersService) generateTokens(managerId, role, companyId string) (*TokensResponse, error) {
	accessToken, err := s.jwtMaker.CreateAccessTokenManager(managerId, role, companyId)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtMaker.CreateRefreshToken(managerId)
	if err != nil {
		return nil, err
	}

	return &TokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthManagersService) Login(req *LoginRequest) (*AuthResponse, error) {
	manager, err := s.managerRepo.FindByLogin(req.Login)
	if err != nil {
		if errors.Is(err, utils.ErrRecordNotFound) {
			return nil, utils.ErrInvalidCredentials
		}
		return nil, err
	}

	if manager == nil {
		return nil, utils.ErrInvalidCredentials
	}

	match, err := utils.Verify(req.Password, manager.Password)

	if err != nil {
		return nil, err
	}

	if !match {
		return nil, utils.ErrInvalidCredentials
	}

	tokens, err := s.generateTokens(manager.ID.String(), manager.Role, manager.CompanyID.String())

	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		Manager:      *s.ToManagerResponse(manager),
	}, nil
}

func (s *AuthManagersService) ChangePassword(oldPassword, newPassword, managerId, ownerId, companyId string) (*ManagerResponse, error) {
	var resp *ManagerResponse
	err := s.db.Transaction(func(tx *gorm.DB) error {
		_, err := s.companyRepo.FindById(ownerId, companyId)
		if err != nil {
			return err
		}

		oldManager, err := s.managerRepo.FindById(managerId)

		if err != nil {
			return err
		}

		if oldManager.CompanyID.String() != companyId {
			return utils.ErrForbidden
		}

		match, err := utils.Verify(oldPassword, oldManager.Password)

		if err != nil {
			return err
		}

		if !match {
			return utils.ErrInvalidCredentials
		}

		newHash, err := utils.Hash(newPassword)

		if err != nil {
			return err
		}

		err = s.managerRepo.UpdatePassword(tx, oldManager, newHash)

		if err != nil {
			return err
		}

		resp = s.ToManagerResponse(oldManager)

		return nil
	})

	return resp, err
}

func (s *AuthManagersService) ToManagerResponse(manager *manager.Manager) *ManagerResponse {
	return &ManagerResponse{
		Login:      manager.Login,
		Name:       manager.Name,
		MiddleName: manager.MiddleName,
		LastName:   manager.LastName,
		CreatedAt:  manager.CreatedAt.Format(time.RFC3339),
	}
}

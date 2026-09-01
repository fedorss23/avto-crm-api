package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	secret []byte
	issuer string
	accessTokenDuration time.Duration
	refreshTokenDuration time.Duration
}

type AccessTokenClaims struct {
	UserID string `json:"user_id"`
	Email string `json:"email"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type PasswordResetClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJWTMaker(secret string, issuer string, accessTokenDuration time.Duration, refreshTokenDuration time.Duration) *JWTMaker {
	return &JWTMaker{
		secret: []byte(secret),
		issuer: issuer,
		accessTokenDuration: accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (m *JWTMaker) CreateAccessToken(userID, email string) (string, error) {
	now := time.Now()
	claims := &AccessTokenClaims{
		UserID: userID,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: m.issuer,
			Subject: userID,
			IssuedAt: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTokenDuration)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *JWTMaker) CreateRefreshToken(userID string) (string, error) {
	now := time.Now()
	claims := &RefreshTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: m.issuer,
			Subject: userID,
			IssuedAt: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.refreshTokenDuration)),
			NotBefore: jwt.NewNumericDate(now),
		},	
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *JWTMaker) CreatePassowrdResestToken(userID string) (string, error) {
	now := time.Now()
	claims := &PasswordResetClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: m.issuer,
			Subject: userID,
			IssuedAt: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Hour)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *JWTMaker) ValidateAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (m *JWTMaker) ValidateRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return m.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*RefreshTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (m *JWTMaker) ValidatePasswordResetToken(tokenString string) (*PasswordResetClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &PasswordResetClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return m.secret, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*PasswordResetClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, fmt.Errorf("invalid token")
}

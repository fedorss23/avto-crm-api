package auth

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
    ErrInvalidCredentials   = errors.New("invalid email or password")
    ErrEmailAlreadyExists   = errors.New("email already exists")
    ErrInvalidToken         = errors.New("invalid token")
    ErrTokenExpired         = errors.New("token expired")
    ErrUserLocked           = errors.New("user account is locked")
    ErrUserInactive         = errors.New("user account is inactive")
    ErrInvalidRefreshToken  = errors.New("invalid refresh token")
    ErrPasswordTooWeak      = errors.New("password too weak")
    ErrOldPasswordIncorrect = errors.New("old password is incorrect")
)

func ErrorToHTTPStatus(err error) int {
    switch err {
    case ErrUserNotFound, ErrInvalidCredentials:
        return 401
    case ErrEmailAlreadyExists:
        return 409
    case ErrInvalidToken, ErrInvalidRefreshToken:
        return 401
    case ErrTokenExpired:
        return 401
    case ErrUserLocked:
        return 423
    case ErrUserInactive:
        return 403
    default:
        return 500
    }
}
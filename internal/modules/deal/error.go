package deal

import (
	"errors"
	"net/http"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrNotAuth            = errors.New("not auth")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidStage       = errors.New("invalid stage")
	ErrRecordNotFound = errors.New("deal not found")
	ErrLastStage = errors.New("stage already last")
	ErrEmptyPipeline = errors.New("empty pipeline")
	ErrInternalError = errors.New("server internal error")
	ErrNotYourDeal = errors.New("you can delete only your deal")
	ErrNotPipeline = errors.New("this deal doesn't have a pipeline")
	ErrNotStages = errors.New("this deal doesn't have stages")
)

func ErrorToHTTPStatus(err error) int {
	switch err {
	case ErrUserNotFound, ErrInvalidCredentials, ErrLastStage, ErrInvalidStage, ErrNotStages, ErrNotPipeline:
		return http.StatusBadRequest
	case ErrNotAuth:
		return http.StatusForbidden
	case ErrRecordNotFound:
		return http.StatusNotFound

	// case ErrEmailAlreadyExists:
	// 	return 409
	// case ErrInvalidToken, ErrInvalidRefreshToken:
	// 	return 401
	// case ErrTokenExpired:
	// 	return 401
	// case ErrUserLocked:
	// 	return 423
	// case ErrUserInactive:
	// 	return 403
	default:
		return http.StatusInternalServerError
	}
}

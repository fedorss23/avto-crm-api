package client

import (
	"errors"
	"net/http"
)

var (
	ErrClientNotFound     = errors.New("client not found")
	ErrNotAuth            = errors.New("not auth")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInternalError      = errors.New("server internal error")
	ErrWithParseOwnerId   = errors.New("error with parse ownerId")
)

var (
	clientNotFoundCode      = "CLIENT_NOT_FOUND"
	notAuthorizedCode       = "NOT_AUTHORIZED"
	defaultCode             = "UNKNOW_ERROR"
	notFoundClientParamCode = "CLIENTID_PARAM_NOT_FOUND"
	pageErrorCode           = "QUERY_PARAM_PAGE_ERROR"
	limitErrorCode          = "QUERY_PARAM_LIMIT_ERROR"
	isFullErrorCode         = "QUERY_PARAM_ISFULL_ERROR"
)

func CodeByError(err error) string {
	switch err {
	case ErrClientNotFound:
		return clientNotFoundCode
	case ErrNotAuth:
		return notAuthorizedCode
	default:
		return defaultCode
	}
}

func ErrorToHTTPStatus(err error) int {
	switch err {
	case ErrInvalidCredentials, ErrWithParseOwnerId:
		return http.StatusBadRequest
	case ErrNotAuth:
		return 401
	case ErrClientNotFound:
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

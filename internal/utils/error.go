package utils

import (
	"errors"
	"fmt"
	"net/http"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrForbidden      = errors.New("access denied")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrConflict       = errors.New("conflict")

	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserLocked         = errors.New("user locked")

	ErrNotPipeline    = errors.New("deal doesn't have pipeline")
	ErrEmptyPipeline  = errors.New("pipeline is empty")
	ErrNotStages      = errors.New("this deal doesn't have stages")
	ErrNotFoundStage  = errors.New("stage not found")
	ErrInvalidStageId = errors.New("invalid stage")
	ErrLastStage      = errors.New("deal has last stage")

	ErrEmailAlredyExists         = errors.New("email already exists")
	ErrInvalidRefreshToken       = errors.New("invalid refresh token")
	ErrPasswordIncorrectLogin    = errors.New("invalid password")
	ErrPasswordIncorrectRegister = errors.New("invalid password")
)

const (
	PageErrorCode   = "QUERY_PARAM_PAGE_ERROR"
	LimitErrorCode  = "QUERY_PARAM_LIMIT_ERROR"
	IsFullErrorCode = "QUERY_PARAM_ISFULL_ERROR"

	ForbiddenCode       = "FORBIDDEN"
	UnauthorizedCode    = "UNAUTHORIZED"
	RecordNotFoundCode  = "RECORD_NOT_FOUND"
	ValidationErrorCode = "VALIDATION_ERROR"
	ConflictCode        = "CONFLICT"
	CredentialsCode     = "INVALID_CREDENTIALS"
	UserLockedCode      = "USER_LOCKED"

	DealPipelineNotFoundCode    = "DEAL_PIPELINE_NOT_FOUND"
	DealPipelineEmptyCode       = "PIPELINE_EMPTY"
	DealStagesNotConfiguredCode = "DEAL_STAGES_NOT_CONFIGURED"
	StageNotFoundCode           = "STAGE_NOT_FOUND"
	InvalidStageCode            = "INVALID_STAGE"
	LastStageCode               = "DEAL_ALREADY_AT_LAST_STAGE"

	InternalErrorCode = "INTERNAL_ERROR"
)

func CodeByError(err error) string {
	switch {
	case errors.Is(err, ErrRecordNotFound):
		return RecordNotFoundCode

	case errors.Is(err, ErrForbidden):
		return ForbiddenCode

	case errors.Is(err, ErrEmailAlredyExists), errors.Is(err, ErrConflict):
		return ConflictCode

	case errors.Is(err, ErrUnauthorized):
		return UnauthorizedCode

	case errors.Is(err, ErrUserNotFound):
		return RecordNotFoundCode

	case errors.Is(err, ErrInvalidCredentials), errors.Is(err, ErrInvalidRefreshToken):
		return UnauthorizedCode

	case errors.Is(err, ErrNotPipeline):
		return DealPipelineNotFoundCode

	case errors.Is(err, ErrEmptyPipeline):
		return DealPipelineEmptyCode

	case errors.Is(err, ErrNotStages):
		return DealStagesNotConfiguredCode

	case errors.Is(err, ErrNotFoundStage):
		return StageNotFoundCode

	case errors.Is(err, ErrInvalidStageId):
		return InvalidStageCode

	case errors.Is(err, ErrLastStage):
		return LastStageCode

	case errors.Is(err, ErrPasswordIncorrectLogin), errors.Is(err, ErrPasswordIncorrectRegister):
		return CredentialsCode

	case errors.Is(err, ErrUserLocked):
		return UserLockedCode

	default:
		return InternalErrorCode
	}
}

func ErrorToHTTPStatus(err error) int {
	switch {
	case errors.Is(err, ErrUnauthorized),
		errors.Is(err, ErrInvalidCredentials),
		errors.Is(err, ErrInvalidRefreshToken):
		return http.StatusUnauthorized

	case errors.Is(err, ErrForbidden),
		errors.Is(err, ErrUserLocked):
		return http.StatusForbidden

	case errors.Is(err, ErrRecordNotFound),
		errors.Is(err, ErrUserNotFound),
		errors.Is(err, ErrNotFoundStage):
		return http.StatusNotFound

	case errors.Is(err, ErrInvalidStageId),
		errors.Is(err, ErrPasswordIncorrectRegister),
		errors.Is(err, ErrEmailAlredyExists),
		errors.Is(err, ErrPasswordIncorrectLogin):
		return http.StatusBadRequest

	case errors.Is(err, ErrNotPipeline),
		errors.Is(err, ErrEmptyPipeline),
		errors.Is(err, ErrNotStages),
		errors.Is(err, ErrLastStage):
		return http.StatusConflict

	default:
		return http.StatusInternalServerError
	}
}

func ParseValidationErrors(err error) map[string]string {
	errs := make(map[string]string)

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			errs[fe.Field()] = validationMessage(fe)
		}
		return errs
	}

	errs["_"] = err.Error()
	return errs
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "field required"
	case "email":
		return "incorrect email"
	case "min":
		return fmt.Sprintf("mininal length/value: %s", fe.Param())
	case "max":
		return fmt.Sprintf("maximal length/value: %s", fe.Param())
	default:
		return fmt.Sprintf("invalid value (%s)", fe.Tag())
	}
}

func ToQueryErrorMessageWithType(name string, queryType string) string {
	return fmt.Sprintf("Error with query-param %s: %s must be %s", name, name, queryType)
}

func ToQueryErrorWithType(name string, queryType string) error {
	return errors.New(fmt.Sprintf("%s must be %s", name, queryType))
}

func ToQueryErrorMessage(name string) string {
	return fmt.Sprintf("Query param %s must be provided", name)
}

func ToQueryError(name string) error {
	return errors.New(fmt.Sprintf("Missing %s", name))
}

func ToParamError(name string) error {
	return errors.New(fmt.Sprintf("Missing %s", name))
}

func ToMissingParamMessage(name string) string {
	return fmt.Sprintf("Param %s must be valid", name)
}

func ValidatePassword(password string) bool {
	if len([]rune(password)) < 8 {
		return false
	}

	var hasUpper, hasLower, hasDigit bool

	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}

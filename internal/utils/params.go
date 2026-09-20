package utils

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetOwnerId(c *gin.Context, ownerId *string) error {
	id := c.GetString("userId")
	if id == "" {
		ErrorResponse(c, 401, "not authorized", errors.New("not authorized"), "NOT_AUTHORIZED")
		return errors.New("not authorized")
	}

	*ownerId = id

	return nil
}

func GetParam(c *gin.Context, name string, cont *string) error {
	param := c.Param(name)

	if param == "" {
		errs := make(map[string]string)
		errs[name] = "param required"
		ValidationErrorResponse(c, errs)
		return ToParamError(name)
	}

	*cont = param

	return nil
}

func GetNumberQuery(c *gin.Context, name string, cont *int, defaultval int, code string) error {
	param, exists := c.GetQuery(name)

	if !exists {
		*cont = defaultval
	} else {
		a, err := strconv.Atoi(param)
		if err != nil {
			ErrorResponse(
				c,
				http.StatusBadRequest,
				ToQueryErrorMessageWithType(name, "number"),
				ToQueryErrorWithType(name, "number"),
				code,
			)
			return ToQueryErrorWithType(name, "number")
		}
		*cont = a
	}

	return nil
}

func GetBoolQuery(c *gin.Context, name string, cont *bool, defaultval bool, code string) error {
	param, exists := c.GetQuery(name)

	if !exists {
		*cont = false
	} else {
		a, err := strconv.ParseBool(param)
		if err != nil {
			ErrorResponse(
				c,
				http.StatusBadRequest,
				ToQueryErrorMessageWithType(name, "true or false"),
				ToQueryErrorWithType(name, "bool"),
				code,
			)
			return ToQueryErrorWithType(name, "bool")
		}
		*cont = a
	}

	return nil
}

func GetStringRequiredQuery(c *gin.Context, name string, cont *string) error {
	param, exists := c.GetQuery(name)
	if !exists {
		var err map[string]string
		err[name] = "query-param required"
		ValidationErrorResponse(
			c,
			err,
		)
		return ToQueryErrorWithType(name, "string")
	}

	*cont = param

	return nil
}

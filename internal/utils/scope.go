package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthScope struct {
	CompanyId uuid.UUID
	ManagerId *uuid.UUID
}

func GetAuthScopeWithManager(c *gin.Context) (*AuthScope, error) {
	var scope AuthScope

	var uid string
	userType, err := GetUserId(c, &uid)
	if err != nil {
		return nil, err
	}

	userId, err := uuid.Parse(uid)

	if err != nil {
		ValidationErrorResponse(c, map[string]string{
			"userId": "invalid id",
		})
		return nil, err
	}

	cp := c.GetString("companyId")

	if cp != "" {
		companyId, err := uuid.Parse(cp)

		if err != nil {
			ValidationErrorResponse(c, map[string]string{
				"companyId": "invalid id",
			})
			return nil, err
		}

		scope.CompanyId = companyId
	}

	if userType == "manager" {
		scope.ManagerId = &userId
	} else {
		manager, ok := c.GetQuery("managerId")

		if ok && manager != "" {
			managerId, err := uuid.Parse(manager)

			if err != nil {
				ValidationErrorResponse(c, map[string]string{
					"managerId": "invalid id",
				})
				return nil, err
			}

			scope.ManagerId = &managerId
		}
	}

	return &scope, nil
}

func GetAuthScope(c *gin.Context) (*AuthScope, error) {
	var scope AuthScope

	var uid string
	userType, err := GetUserId(c, &uid)
	if err != nil {
		return nil, err
	}

	userId, err := uuid.Parse(uid)

	if err != nil {
		ValidationErrorResponse(c, map[string]string{
			"userId": "invalid id",
		})
		return nil, err
	}

	cp := c.GetString("companyId")

	if cp != "" {
		companyId, err := uuid.Parse(cp)

		if err != nil {
			ValidationErrorResponse(c, map[string]string{
				"companyId": "invalid id",
			})
			return nil, err
		}

		scope.CompanyId = companyId
	}

	if userType == "manager" {
		scope.ManagerId = &userId
	}

	return &scope, nil
}

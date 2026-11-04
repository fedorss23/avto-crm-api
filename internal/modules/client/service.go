package client

import (
	"avto-crm-api/internal/utils"

	"gorm.io/gorm"
)

type ClientService struct {
	clientRepo *ClientRepository
	db *gorm.DB
}

func NewClientService(clientRepo *ClientRepository, db *gorm.DB) *ClientService {
	return &ClientService{
		clientRepo: clientRepo,
		db: db,
	}
}

func (s *ClientService) FindList(page, limit int, scope *utils.AuthScope) ([]Client, int64, error) {
	return s.clientRepo.FindList(page, limit, scope)
}

func (s *ClientService) Update(req *UpdateClientRequest, id string, scope *utils.AuthScope) (*Client, error) {
	var client *Client
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		client, err = s.clientRepo.FindForUpdate(tx, id, scope)

		if err != nil {
			return err
		}

		updateMap := make(map[string]interface{})

		if req.Email != nil {
			updateMap["email"] = *req.Email
		}

		if req.Phone != nil {
			updateMap["phone"] = *req.Phone
		}

		if len(updateMap) == 0 {
			return utils.ErrEmptyBodyForUpdate
		}

		return tx.Model(client).Updates(updateMap).Error
	})

	return client, err
}

func (s *ClientService) Delete(id string, scope *utils.AuthScope) error {
	return s.clientRepo.Delete(id, scope)
}
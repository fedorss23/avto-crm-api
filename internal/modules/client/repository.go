package client

import "gorm.io/gorm"

type ClientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) *ClientRepository {
	return &ClientRepository{
		db: db,
	}
}

func (r *ClientRepository) Create(req *Client) error {
	return r.db.Create(req).Error
}

func (r *ClientRepository) FindByOwnerId(ownerId string) ([]Client, error) {
	var clients []Client

	if err := r.db.Where("ownerId = ?", ownerId).Find(&clients).Error; err != nil {
		return nil, err
	}

	return clients, nil
}

func (r *ClientRepository) Update(req *Client) error {
	return r.db.Save(req).Error
}

func (r *ClientRepository) FindById(clientId string) (*Client, error) {
	var client *Client

	if err := r.db.Where("id = ?", clientId).First(&client).Error; err != nil {
		return nil, err
	}

	return client, nil
}

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

func (r *ClientRepository) CreateWithTransaction(tx *gorm.DB, req *Client) error {
	return tx.Create(req).Error
}

func (r *ClientRepository) FindByOwnerId(ownerId string, page, limit int) ([]Client, int64, error) {
	var clients []Client
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Model(&Client{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&clients).Error

	return clients, total, err
}

func (r *ClientRepository) Update(req *Client) error {
	return r.db.Save(req).Error
}

func (r *ClientRepository) FindById(clientId, ownerId string) (*Client, error) {
	var client *Client

	if err := r.db.Where("id = ? AND owner_id = ?", clientId, ownerId).First(&client).Error; err != nil {
		return nil, err
	}

	return client, nil
}

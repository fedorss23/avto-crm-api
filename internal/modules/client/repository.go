package client

import (
	"avto-crm-api/internal/utils"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ClientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) *ClientRepository {
	return &ClientRepository{
		db: db,
	}
}

func (r *ClientRepository) CreateWithTransaction(tx *gorm.DB, req *Client) error {
	return tx.Create(req).Error
}

func (r *ClientRepository) FindList(page, limit int, scope *utils.AuthScope) ([]Client, int64, error) {
	var clients []Client
	var total int64

	offset := (page - 1) * limit

	query := r.db.Model(&Client{})

	if scope.CompanyId != nil {
		query = query.Where("company_id = ?", *scope.CompanyId)
	}

	if scope.ManagerId != nil {
		query = query.Where("manager_id = ?", *scope.ManagerId)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&clients).Error

	return clients, total, err
}

func (r *ClientRepository) FindForUpdate(tx *gorm.DB, id string, scope *utils.AuthScope) (*Client, error) {
	var client Client

	query := tx.Clauses(clause.Locking{
		Strength: "UPDATE",
	}).Where("id = ?", id)

	if scope.CompanyId != nil {
		query = query.Where("company_id = ?", scope.CompanyId)
	}

	if scope.ManagerId != nil {
		query = query.Where("manager_id = ?", scope.ManagerId)
	}

	err := query.First(&client).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrRecordNotFound
		}

		return nil, err
	}

	return &client, nil
}

func (r *ClientRepository) FindById(id string, scope *utils.AuthScope) (*Client, error) {
	var client Client

	query := r.db.Where("id = ?", id)

	if scope.CompanyId != nil {
		query = query.Where("company_id = ?", *scope.CompanyId)
	}

	if scope.ManagerId != nil {
		query = query.Where("manager_id = ?", *scope.ManagerId)
	}

	err := query.First(&client).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrRecordNotFound
		}

		return nil, err
	}

	return &client, nil
}

func (r *ClientRepository) Delete(tx *gorm.DB, id uuid.UUID, scope *utils.AuthScope) error {
	query := tx.Where("id = ?", id)

	if scope.CompanyId != nil {
		query = query.Where("company_id = ?", scope.CompanyId)
	}

	if scope.ManagerId != nil {
		query = query.Where("manager_id = ?", scope.ManagerId)
	}

	result := query.Delete(&Client{})

	if result.RowsAffected == 0 {
		return utils.ErrRecordNotFound
	}

	if result.Error != nil {
		return result.Error
	}

	return nil
}
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

func (r *ClientRepository) Create(tx *gorm.DB, req *Client) error {
	return tx.Create(req).Error
}

func (r *ClientRepository) FindList(page, limit int, scope *utils.AuthScope) ([]Client, int64, error) {
	var clients []Client
	var total int64

	offset := (page - 1) * limit

	query := r.db.Model(&Client{}).
		Joins("JOIN deals ON deals.id = clients.deal_id").
		Where("deals.company_id = ?", scope.CompanyId)

	if scope.ManagerId != nil {
		query = query.Where("deals.manager_id = ?", *scope.ManagerId)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&clients).Error

	return clients, total, err
}

func (r *ClientRepository) FindForUpdate(tx *gorm.DB, id string, scope *utils.AuthScope) (*Client, error) {
	var client Client

	err := tx.Clauses(clause.Locking{
		Strength: "UPDATE",
	}).
		Where("id = ?", id).
		Where("EXISTS (?)", r.checkAccessQuery(tx, scope)).
		First(&client).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrRecordNotFound
		}

		return nil, err
	}

	return &client, nil
}

func (r *ClientRepository) DeleteWithDeal(tx *gorm.DB, id uuid.UUID) error {
	query := tx.Where("id = ?", id)

	result := query.Delete(&Client{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *ClientRepository) Delete(id string, scope *utils.AuthScope) error {
	result := r.db.
		Where("id = ?", id).
		Where("EXISTS (?)", r.checkAccessQuery(r.db, scope)).
		Unscoped().
		Delete(&Client{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return utils.ErrRecordNotFound
	}

	return nil
}

func (r *ClientRepository) checkAccessQuery(tx *gorm.DB, scope *utils.AuthScope) *gorm.DB {
	subQuery := tx.Table("deals").
		Select("1").
		Where("deals.id = clients.deal_id").
		Where("deals.company_id = ?", scope.CompanyId)

	if scope.ManagerId != nil {
		subQuery = subQuery.Where("deals.manager_id = ?", *scope.ManagerId)
	}
	return subQuery
}

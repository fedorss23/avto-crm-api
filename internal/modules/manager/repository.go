package manager

import (
	"avto-crm-api/internal/utils"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ManagerRepository struct {
	db *gorm.DB
}

func NewManagerRepository(db *gorm.DB) *ManagerRepository {
	return &ManagerRepository{
		db: db,
	}
}

func (r *ManagerRepository) FindById(managerId string) (*Manager, error) {
	var manager *Manager
	if err := r.db.Where("id = ?", managerId).First(manager).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrRecordNotFound
		}

		return nil, err
	}

	return manager, nil
}

func (r *ManagerRepository) Create(manager *Manager) error {
	return r.db.Create(manager).Error
}

func (r *ManagerRepository) FindByLogin(login string) (*Manager, error) {
	var manager Manager
	if err := r.db.Where("login = ?", login).First(&manager).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrRecordNotFound
		}

		return nil, err
	}

	return &manager, nil
}

func (r *ManagerRepository) UpdatePassword(tx *gorm.DB, manager *Manager, hash string) error {
	return tx.Model(manager).Update("password", hash).Error
}

func (r *ManagerRepository) GetFreeManager(tx *gorm.DB, companyId uuid.UUID) (string, error) {
	if err := r.CheckManagers(tx, companyId); err != nil {
		return "", err
	}

	var manager Manager
	err := tx.Model(&Manager{}).
		Joins("LEFT JOIN deals ON deals.manager_id = managers.id").
		Where("managers.company_id = ?::uuid", companyId).
		Group("managers.id").
		Order("COUNT(deals.id) DESC",).
		First(&manager).Error

	if err != nil {
		return "", err
	}

	log.Println(manager.ID)

	return manager.ID.String(), nil
}

func (r *ManagerRepository) CheckManagers(tx *gorm.DB, companyId uuid.UUID) error {
	var manager Manager
	if err := tx.Model(&Manager{}).Where("company_id = ?::uuid", companyId).First(&manager).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrNotManagers
		}
		return err
	}
	fmt.Println(manager.ID, manager.Name)
	return nil
}

package user

import (
	"errors"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}
	
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FindById(id string) (*User, error) {
	var user *User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindByEmail(email string) (*User, error) {
	var user *User
	if err := r.db.Where("email = ? ", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Create(user *User) error {
	return r.db.Create(user).Error
} 

func (r *UserRepository) Update(user *User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id string) error {
	return r.db.Delete(&User{}, "id = ?", id).Error
}

func (r *UserRepository) FindList(page, limit int) ([]User, int64, error) {
	var users []User
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Model(&User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error

	return users, total, err
}
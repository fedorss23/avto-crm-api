package car

import "gorm.io/gorm"

type CarRepository struct {
	db *gorm.DB
}

func NewCarRepository(db *gorm.DB) *CarRepository {
	return &CarRepository{
		db: db,
	}
}

func (r *CarRepository) Create(car *Car) error {
	return r.db.Create(car).Error
}

func (r *CarRepository) CreateWithTx(tx *gorm.DB, car *Car) error {
	return tx.Create(car).Error
}

func (r *CarRepository) CreateInTransaction(tx *gorm.DB, car *Car) error {
	return tx.Create(car).Error
}

func (r *CarRepository) FindList(page, limit int) ([]Car, int64, error) {
	var cars []Car
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Model(&Car{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&cars).Error

	return cars, total, err
}

func (r *CarRepository) Update(car *Car) error {
	return r.db.Save(car).Error
}

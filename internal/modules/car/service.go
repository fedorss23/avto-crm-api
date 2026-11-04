package car

import (
	"avto-crm-api/internal/utils"

	"gorm.io/gorm"
)

type CarService struct {
	carRepo *CarRepository
	db *gorm.DB
}

func NewCarService(carRepo *CarRepository, db *gorm.DB) *CarService {
	return &CarService{
		carRepo: carRepo,
		db: db,
	}
}

func (s *CarService) FindList(page, limit int, scope *utils.AuthScope) ([]Car, int64, error) {
	return s.carRepo.FindList(page, limit, scope)
}

func (s *CarService) Update(req *UpdateCarRequest, id string, scope *utils.AuthScope) (*Car, error) {
	var car *Car
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		car, err = s.carRepo.FindForUpdate(tx, id, scope)
		if err != nil {
			return err
		}

		updateMap := make(map[string]interface{})

		if req.Model != nil {
			updateMap["model"] = req.Model
		}

		if len(updateMap) == 0 {
			return utils.ErrEmptyBodyForUpdate
		}

		return tx.Model(car).Updates(updateMap).Error
	})

	return car, err
}
func (s *CarService) Delete(id string, scope *utils.AuthScope) error {
	return s.carRepo.Delete(id, scope)
}
package deal

import (
	"avto-crm-api/internal/modules/car"
	"avto-crm-api/internal/modules/pipeline"
	"avto-crm-api/internal/modules/stage"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DealService struct {
	db        *gorm.DB
	dealRepo  *DealRepository
	pipeRepo  *pipeline.PipelineRepository
	stageRepo *stage.StageRepository
	carRepo   *car.CarRepository
}

func NewDealService(db *gorm.DB, dealRepo *DealRepository, carRepo *car.CarRepository) *DealService {
	return &DealService{
		db:       db,
		dealRepo: dealRepo,
		carRepo:  carRepo,
	}
}

func (s *DealService) SetNextPage(ownerID, dealID string) (*Deal, error) {
	return s.dealRepo.SetNextStage(s.db, ownerID, dealID)
}

func (s *DealService) FindAll(page, limit int, isFull bool) ([]Deal, int64, error) {
	if isFull {
		return s.dealRepo.FindListWithAll(s.db, page, limit)
	} else {
		return s.dealRepo.FindList(s.db, page, limit)
	}
}

func (s *DealService) CreateFullDeal(req *CreateDealRequest, ownerId string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		pipeline := &pipeline.Pipeline{
			Name:        req.Pipeline.Name,
			Source:      req.Pipeline.Source,
			Destination: req.Pipeline.Destination,
		}

		var firstStageID uuid.UUID

		car := &car.Car{
			Model: req.Car.Model,
		}

		oid, err := uuid.Parse(ownerId)

		if err != nil {
			return err
		}

		deal := &Deal{
			Name:         req.Name,
			Pipeline:     pipeline,
			OwnerID:      oid,
			CurrentStage: &firstStageID,
			Car:          car,
		}

		if err := s.dealRepo.Create(tx, deal); err != nil {
			return err
		}

		for index, i := range req.Pipeline.Stages {
			stage := stage.Stage{
				Name:       i.Name,
				PipelineID: pipeline.ID,
			}

			if i.Description != "" {
				stage.Description = &i.Description
			}

			if err := s.stageRepo.Create(tx, &stage); err != nil {
				return err
			}

			if index == 0 {
				firstStageID = stage.ID
			}
 		}

		return nil
	})
}

func (s *DealService) Update(req *Deal) error {
	return s.dealRepo.Update(s.db, req)
}

func (s *DealService) FindDealByOwnerId(ownerID string) ([]Deal, int64, error) {
	return s.dealRepo.FindByOwnerId(s.db, ownerID)
}

func (s *DealService) FindDealByClientId(clientID string) ([]Deal, int64, error) {
	return s.dealRepo.FindByClientID(s.db, clientID)
}

func (s *DealService) Delete(ownerID, dealID string) error {
	return s.dealRepo.Delete(s.db, ownerID, dealID)
}

func (s *DealService) FindById(id string, isFull bool) (*Deal, error) {
	if isFull {
		return s.dealRepo.FindFullById(s.db, id)
	} else {
		return s.dealRepo.FindById(s.db, id)
	}
}
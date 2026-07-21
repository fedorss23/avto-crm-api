package deal

import (
	"avto-crm-api/internal/modules/pipeline"
	"avto-crm-api/internal/modules/stage"

	"gorm.io/gorm"
)

type DealService struct {
	db *gorm.DB
	dealRepo *DealRepository
	pipeRepo *pipeline.PipelineRepository
	stageRepo *stage.StageRepository
}

func NewDealService(db *gorm.DB, dealRepo *DealRepository) *DealService {
	return &DealService{
		db: db,
		dealRepo: dealRepo,
	}
}

func (s *DealService) CreateFullDeal(req *CreateDealRequest) (error) {
	return s.db.Transaction(func(tx *gorm.DB) error {
		pipeline := &pipeline.Pipeline{
			Name: req.PipelineName,
			Source: req.Source,
			Destination: req.Destination,
		}

		err := s.pipeRepo.Create(tx, pipeline)
		if err != nil {
			return err
		}

		for _, i := range req.Stages {
			stage := stage.Stage{
				Name: i.Name,
				PipelineID: pipeline.ID,
			}

			if i.Description != "" {
				stage.Description = &i .Description
			}

			if err := s.stageRepo.Create(tx, &stage); err != nil {
				return err
			}
		}

		deal := &Deal{
			Name: req.Name,
			PipelineID: &pipeline.ID,
			OwnerID: req.OwnerID,
		}

		if err := s.dealRepo.Create(tx, deal); err != nil {
			return  err
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
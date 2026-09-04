package deal

import (
	"avto-crm-api/internal/modules/car"
	"avto-crm-api/internal/modules/client"
	"avto-crm-api/internal/modules/pipeline"
	"avto-crm-api/internal/modules/stage"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DealService struct {
	db        *gorm.DB
	dealRepo  *DealRepository
	pipeRepo  *pipeline.PipelineRepository
	stageRepo *stage.StageRepository
	carRepo   *car.CarRepository
	clientRepo *client.ClientRepository
}

func NewDealService(
	db *gorm.DB, 
	dealRepo *DealRepository, 
	carRepo *car.CarRepository, 
	pipeRepo *pipeline.PipelineRepository, 
	stageRepo *stage.StageRepository,
	clientRepo *client.ClientRepository,
) *DealService {
	return &DealService{
		db:       db,
		dealRepo: dealRepo,
		carRepo:  carRepo,
		pipeRepo: pipeRepo,
		stageRepo: stageRepo,
		clientRepo: clientRepo,
	}
}

func (s *DealService) SetNextPage(ownerID, dealID string) (*Deal, error) {
	return s.dealRepo.SetNextStage(ownerID, dealID)
}

func (s *DealService) FindAll(page, limit int, isFull bool) ([]Deal, int64, error) {
	if isFull {
		return s.dealRepo.FindListWithAll(page, limit)
	} else {
		return s.dealRepo.FindList(page, limit)
	}
}

func (s *DealService) CreateFullDeal(req *CreateDealRequest, ownerId string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		pipeline := &pipeline.Pipeline{
			Name:        req.Pipeline.Name,
			Source:      req.Pipeline.Source,
			Destination: req.Pipeline.Destination,
		}

		car := &car.Car{
			Model: req.Car.Model,
		}

		oid, err := uuid.Parse(ownerId)

		if err != nil {
			return err
		}

		client := &client.Client{
			Name: req.Client.Name,
			OwnerID: oid,
		}

		if req.Client.Email != nil  {
			client.Email = req.Client.Email
		}

		if req.Client.Phone != nil {
			client.Phone = req.Client.Phone
		}

		now := time.Now()

		dueDate := now.AddDate(0, 0, req.Term)

		deal := &Deal{
			Name:     req.Name,
			Pipeline: pipeline,
			OwnerID:  oid,
			Car:      car,
			Client: client,
			Term: req.Term,
			DueDate: dueDate,
			Total: req.Total,
		}

		if err := s.dealRepo.Create(tx, deal); err != nil {
			return err
		}

		var firstStageID uuid.UUID
		var currentStageName string

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
				currentStageName = stage.Name
			}
		}

		deal.PipelineId = &pipeline.ID
		deal.ClientId = &client.ID
		deal.CarId = &car.ID
		deal.CurrentStageId = &firstStageID
		deal.CurrentStageName = &currentStageName

		if err := s.dealRepo.Update(tx, deal); err != nil {
			return err
		}

		return nil
	})
}

func (s *DealService) Update(req *Deal) error {
	return s.dealRepo.Update(s.db, req)
}

func (s *DealService) FindDealByOwnerId(ownerID string) ([]Deal, int64, error) {
	return s.dealRepo.FindByOwnerId(ownerID)
}

func (s *DealService) FindDealByClientId(clientID string) ([]Deal, int64, error) {
	return s.dealRepo.FindByClientID(clientID)
}

func (s *DealService) Delete(ownerID, dealID string) error {
	return s.dealRepo.Delete(ownerID, dealID)
}

func (s *DealService) FindById(id string, isFull bool) (*Deal, error) {
	if isFull {
		return s.dealRepo.FindFullById(id)
	} else {
		return s.dealRepo.FindById(id)
	}
}

func (s *DealService) ChangeStatus(ownerID, dealID, status string) error {
	return s.dealRepo.ChangeStatus(ownerID, dealID, status)
}

func (s *DealService) ChangeStage(ownerId, dealId, stageId string) error {
	return s.dealRepo.SetStage(ownerId, dealId, stageId)
}
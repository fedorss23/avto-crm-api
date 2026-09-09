package deal

import (
	"avto-crm-api/internal/modules/car"
	"avto-crm-api/internal/modules/client"
	"avto-crm-api/internal/modules/pipeline"
	"avto-crm-api/internal/modules/stage"
	"avto-crm-api/internal/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DealService struct {
	db         *gorm.DB
	dealRepo   *DealRepository
	pipeRepo   *pipeline.PipelineRepository
	stageRepo  *stage.StageRepository
	carRepo    *car.CarRepository
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
		db:         db,
		dealRepo:   dealRepo,
		carRepo:    carRepo,
		pipeRepo:   pipeRepo,
		stageRepo:  stageRepo,
		clientRepo: clientRepo,
	}
}

func (s *DealService) SetNextStage(ownerId, dealId string) (*Deal, error) {
	deal, err := s.dealRepo.FindFullById(dealId, ownerId)

	if err != nil {
		return nil, err
	}

	if deal.Pipeline == nil {
		return nil, utils.ErrNotPipeline
	}

	if deal.Pipeline.Stages == nil {
		return nil, utils.ErrNotStages
	}

	stages := deal.Pipeline.Stages

	if len(stages) == 0 {
		return nil, utils.ErrEmptyPipeline
	}

	if deal.CurrentStageId == nil {
		if err := s.dealRepo.SetNextStage(deal, stages[0].ID.String()); err != nil {
			return nil, err
		}

		deal.CurrentStageId = &stages[0].ID
		deal.CurrentStageName = &stages[0].Name
		return deal, nil
	}

	for i := 0; i < len(stages); i++ {
		if stages[i].ID.String() == deal.CurrentStageId.String() {
			if i == len(stages)-1 {
				return deal, utils.ErrLastStage
			}
			if err :=  s.dealRepo.SetNextStage(deal, stages[i+1].ID.String()); err != nil {
				return nil, err
			}
			deal.CurrentStageId = &stages[i+1].ID
			deal.CurrentStageName = &stages[i+1].Name
			return deal, nil
		}
	}

	return nil, utils.ErrInvalidStageId
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
			Name:    req.Client.Name,
			OwnerID: oid,
		}

		if req.Client.Email != nil {
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
			Client:   client,
			Term:     req.Term,
			DueDate:  dueDate,
			Total:    req.Total,
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

		return s.dealRepo.TransactionUpdate(tx, deal)
	})
}

func (s *DealService) Update(req *Deal, ownerId string) error {
	if req.OwnerID.String() != ownerId {
		return utils.ErrForbidden
	}

	_, err := s.dealRepo.FindById(req.ID.String(), ownerId)
	if err != nil {
		return err
	}

	if req.CurrentStageId == nil {
		return utils.ErrInvalidStageId
	}

	for _, i := range req.Pipeline.Stages {
		if i.ID == *req.CurrentStageId {
			req.CurrentStageName = &i.Name
			return s.dealRepo.Update(req)
		}
	}

	return utils.ErrNotFoundStage
}

func (s *DealService) FindDealByOwnerId(ownerId string) ([]Deal, int64, error) {
	return s.dealRepo.FindByOwnerId(ownerId)
}

func (s *DealService) FindDealByClientId(ownerId, clientId string) ([]Deal, int64, error) {
	return s.dealRepo.FindByClientID(clientId, ownerId)
}

func (s *DealService) Delete(ownerId, dealId string) error {
	return s.dealRepo.Delete(ownerId, dealId)
}

func (s *DealService) FindById(ownerId, dealId string, isFull bool) (*Deal, error) {
	if isFull {
		return s.dealRepo.FindFullById(dealId, ownerId)
	} else {
		return s.dealRepo.FindById(dealId, ownerId)
	}
}

func (s *DealService) ChangeStatus(ownerId, dealId, status string) error {
	deal, err := s.dealRepo.FindById(dealId, ownerId)

	if err != nil {
		return err
	}

	deal.Status = status

	return s.dealRepo.Update(deal)
}

func (s *DealService) ChangeStage(ownerId, dealId, stageId string) error {
	deal, err := s.dealRepo.FindFullById(dealId, ownerId)
	if err != nil {
		return err
	}

	parsedId, err := uuid.Parse(stageId)

	if err != nil {
		return utils.ErrInvalidStageId
	}

	if deal.Pipeline == nil {
		return utils.ErrNotPipeline
	}

	if deal.Pipeline.Stages == nil {
		return utils.ErrNotStages
	}

	stages := deal.Pipeline.Stages

	if len(stages) == 0 {
		return utils.ErrEmptyPipeline
	}

	for _, i := range deal.Pipeline.Stages {
		if i.ID == parsedId {
			deal.CurrentStageId = &parsedId
			deal.CurrentStageName = &i.Name
			return s.dealRepo.Update(deal)
		}
	}

	return utils.ErrNotFoundStage
}

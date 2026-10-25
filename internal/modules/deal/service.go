package deal

import (
	"avto-crm-api/internal/modules/car"
	"avto-crm-api/internal/modules/client"
	"avto-crm-api/internal/modules/manager"
	"avto-crm-api/internal/modules/pipeline"
	"avto-crm-api/internal/modules/stage"
	"avto-crm-api/internal/utils"
	"errors"
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
	managerRepo *manager.ManagerRepository
}

func NewDealService(
	db *gorm.DB,
	dealRepo *DealRepository,
	carRepo *car.CarRepository,
	pipeRepo *pipeline.PipelineRepository,
	stageRepo *stage.StageRepository,
	clientRepo *client.ClientRepository,
	managerRepo *manager.ManagerRepository,
) *DealService {
	return &DealService{
		db:         db,
		dealRepo:   dealRepo,
		carRepo:    carRepo,
		pipeRepo:   pipeRepo,
		stageRepo:  stageRepo,
		clientRepo: clientRepo,
		managerRepo: managerRepo,
	}
}

func (s *DealService) GetTotalByStatus(scope *DealScope, status string) (int64, error) {
	return s.dealRepo.GetTotalForStatus(scope, status)
}

func (s *DealService) CreateFullDeal(req *CreateDealRequest, scope *DealScope) (*Deal, error) {
	var deal *Deal
	err := s.db.Transaction(func(tx *gorm.DB) error {
		cid, err := uuid.Parse(scope.CompanyId)
		if err != nil {
			return err
		}

		if scope.ManagerId == nil {
			id, err := s.managerRepo.GetFreeManager(tx, cid)
			if err != nil {
				if errors.Is(err, utils.ErrRecordNotFound) {
					return utils.ErrNotManagers
				}
				return err
			}
			scope.ManagerId = &id
		}

		mid, err := uuid.Parse(*scope.ManagerId)
		if err != nil {
			return err
		}

		if len(req.Pipeline.Stages) == 0 {
			return utils.ErrNotFoundStage
		}

		newPipeline := &pipeline.Pipeline{
			Name:        req.Pipeline.Name,
			Source:      req.Pipeline.Source,
			Destination: req.Pipeline.Destination,
		}

		for index, st := range req.Pipeline.Stages {
			stg := stage.Stage{
				Name:   st.Name,
				Number: index + 1,
			}
			if st.Description != "" {
				stg.Description = &st.Description
			}
			newPipeline.Stages = append(newPipeline.Stages, stg)
		}

		newCar := &car.Car{Model: req.Car.Model}

		newClient := &client.Client{
			Name:    req.Client.Name,
			OwnerID: mid,
		}

		if req.Client.Email != nil {
			newClient.Email = req.Client.Email
		}

		if req.Client.Phone != nil {
			newClient.Phone = req.Client.Phone
		}

		deal = &Deal{
			Name:      req.Name,
			Pipeline:  newPipeline,
			Client:    newClient,
			Car:       newCar,
			ManagerID: mid,
			CompanyID: cid,
			Term:      req.Term,
			DueDate:   time.Now().AddDate(0, 0, req.Term),
			Total:     req.Total,
		}

		if err := s.dealRepo.Create(tx, deal); err != nil {
			return err
		}

		firstStage := deal.Pipeline.Stages[0]

		return tx.Model(deal).Updates(map[string]interface{}{
			"current_stage_id":   &firstStage.ID,
			"current_stage_name": &firstStage.Name,
		}).Error
	})

	return deal, err
}

func (s *DealService) Update(req *UpdateDealRequest, dealId string, scope *DealScope) (*Deal, error) {
	var deal *Deal
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		deal, err = s.dealRepo.FindFullForUpdate(tx, dealId, scope)

		if err != nil {
			return err
		}

		updateMap := make(map[string]interface{})

		if req.Name != nil {
			updateMap["name"] = *req.Name
		}

		if req.Status != nil {
			updateMap["status"] = *req.Status
		}

		if req.DueDate != nil {
			parseDate, err := time.Parse("01-02-2006", *req.DueDate)
			if err != nil {
				return err
			}
			updateMap["due_date"] = parseDate
		}

		if req.Term != nil {
			updateMap["term"] = *req.Term
		}

		if req.Total != nil {
			updateMap["total"] = *req.Total
		}

		if req.CurrentStageId != nil {
			parseId, err := uuid.Parse(*req.CurrentStageId)
			if err != nil {
				return err
			}

			isFound := false

			for _, i := range deal.Pipeline.Stages {
				if i.ID == parseId {
					stage := i
					updateMap["current_stage_id"] = parseId
					updateMap["current_stage_name"] = stage.Name
					isFound = true
					break
				}
			}

			if !isFound {
				return utils.ErrNotFoundStage
			}
		}

		return tx.Model(deal).Updates(updateMap).Error
	})

	return deal, err
}

func (s *DealService) FindDealByOwnerId(filters *DealFilters) ([]Deal, int64, error) {
	return s.dealRepo.FindList(filters)
}

func (s *DealService) Delete(dealId string, scope *DealScope) error {
	return s.dealRepo.Delete(dealId, scope)
}

func (s *DealService) FindById(dealId string, scope *DealScope, isFull bool) (*Deal, error) {
	return s.dealRepo.FindById(dealId, scope, isFull)
}

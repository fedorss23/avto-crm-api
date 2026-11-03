package pipeline

import (
	"avto-crm-api/internal/modules/stage"
	"avto-crm-api/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PipelineService struct {
	pipeRepo  *PipelineRepository
	stageRepo *stage.StageRepository
	db        *gorm.DB
}

func NewPipelineService(pipeRepo *PipelineRepository, stageRepo *stage.StageRepository, db *gorm.DB) *PipelineService {
	return &PipelineService{
		pipeRepo:  pipeRepo,
		stageRepo: stageRepo,
		db:        db,
	}
}

func (p *PipelineService) Create(req *Pipeline) error {
	return p.pipeRepo.Create(p.db, req)
}

func (p *PipelineService) Update(pipelineId string, req *UpdatePipelineRequest, scope *utils.AuthScope) (*Pipeline, error) {
	var pipeline *Pipeline
	err := p.db.Transaction(func(tx *gorm.DB) error {
		pid, err := uuid.Parse(pipelineId)
		if err != nil {
			return err
		}

		pipeline, err = p.pipeRepo.FindForUpdate(tx, pid, scope, false)
		if err != nil {
			return err
		}

		updateMap := make(map[string]interface{})

		if req.Name != nil {
			updateMap["name"] = *req.Name
		}

		if req.Source != nil {
			updateMap["source"] = *req.Source
		}

		if req.Destination != nil {
			updateMap["destination"] = *req.Destination
		}

		return tx.Model(pipeline).Updates(updateMap).Error
	})
	return pipeline, err
}

func (p *PipelineService) AddStages(pipelineId string, req *AddPipelineStages, scope *utils.AuthScope) (*Pipeline, error) {
	var pipeline *Pipeline
	err := p.db.Transaction(func(tx *gorm.DB) error {
		pid, err := uuid.Parse(pipelineId)
		if err != nil {
			return err
		}

		pipeline, err = p.pipeRepo.FindForUpdate(tx, pid, scope, true)
		if err != nil {
			return err
		}

		lastN := pipeline.Stages[len(pipeline.Stages)-1].Number

		for _, i := range req.Stages {
			st := stage.Stage{
				Name:   i.Name,
				Number: lastN + 1,
			}

			if i.Description != nil {
				st.Description = i.Description
			}

			pipeline.Stages = append(pipeline.Stages, st)
			lastN++
		}

		return tx.Save(pipeline).Error
	})

	return pipeline, err
}

func (p *PipelineService) DeleteStage(pipelineId, stageId string, scope *utils.AuthScope) (*Pipeline, error) {
	var pipeline *Pipeline
	err := p.db.Transaction(func(tx *gorm.DB) error {
		pid, err := uuid.Parse(pipelineId)
		if err != nil {
			return err
		}

		sti, err := uuid.Parse(stageId)
		if err != nil {
			return err
		}

		pipeline, err = p.pipeRepo.FindForUpdate(tx, pid, scope, true)
		if err != nil {
			return err
		}

		in := -1

		for index, i := range pipeline.Stages {
			if in != -1 {
				i.Number--
			} else if i.ID.String() == stageId {
				in = index
			}
		}

		if in == -1 {
			return utils.ErrRecordNotFound
		}

		if in == len(pipeline.Stages)-1 {
			pipeline.Stages = pipeline.Stages[:in]
		} else if in == 0 {
			pipeline.Stages = pipeline.Stages[in:]
		} else {
			pipeline.Stages = append(pipeline.Stages[:in], pipeline.Stages[in+1:]...)
		}

		ferr := tx.Save(pipeline).Error
		if ferr != nil {
			return ferr
		}

		return p.stageRepo.Delete(tx, sti, scope)
	})
	return pipeline, err
}

func (p *PipelineService) FindAll(page, limit int) ([]Pipeline, int64, error) {
	return p.pipeRepo.FindList(p.db, page, limit)
}

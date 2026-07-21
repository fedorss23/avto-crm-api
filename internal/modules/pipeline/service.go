package pipeline

import (
	"gorm.io/gorm"
)

type PipelineService struct {
	pipeRepo *PipelineRepository
	db *gorm.DB
}

func NewPipelineService(pipeRepo *PipelineRepository, db *gorm.DB) *PipelineService {
	return &PipelineService{
		pipeRepo: pipeRepo,
		db: db,
	}
}

func (p *PipelineService) Create(req *Pipeline) error {
	return p.pipeRepo.Create(p.db, req)
}

func (p *PipelineService) Update(req *Pipeline) error {
	return p.pipeRepo.Update(p.db, req)
}

func (p *PipelineService) FindAll(page, limit int) ([]Pipeline, int64, error) {
	return p.pipeRepo.FindList(p.db, page, limit)
}

func (p *PipelineService) FindById(id string) (*Pipeline, error) {
	return p.pipeRepo.FindById(p.db, id)
}
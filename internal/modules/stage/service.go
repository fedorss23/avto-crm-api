package stage

type StageService struct {
	stageRepo *StageRepository
}

func NewStageService(stageRepo *StageRepository) *StageService {
	return &StageService{
		stageRepo: stageRepo,
	}
}

package stage

type StageService struct {
	stageRepo *StageRepository
}

func NewStageService(stageRepo *StageRepository) *StageService {
	return &StageService{
		stageRepo: stageRepo,
	}
}

func (s *StageService) Update(req *UpdateStageRequest, stageId string, ownerId string) (*Stage, error) {
	data := make(map[string]interface{})

	if req.Destination != nil {
		data["destination"] = *req.Destination
	}

	if req.Name != nil {
		data["name"] = *req.Name
	}

	if req.Number != nil {
		data["number"] = *req.Number
	}

	return s.stageRepo.Update(data, stageId, ownerId)
}
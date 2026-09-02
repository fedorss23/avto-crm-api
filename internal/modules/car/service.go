package car 

type CarService struct {
	carRepo *CarRepository
}

func NewCarService(carRepo *CarRepository) *CarService {
	return &CarService{
		carRepo: carRepo,
	}
}

func (s *CarService) Create(req *Car) error {
	return s.carRepo.Create(req)
}

func (s *CarService) FindList(page, limit int) ([]Car, int64, error) {
	return s.carRepo.FindList(page, limit)
}

func (s *CarService) Update(req *Car) error {
	return s.carRepo.Update(req)
}
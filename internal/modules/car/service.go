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
	return s.Create(req)
}

func (s *CarService) FindList(page, limit int) ([]Car, int64, error) {
	return s.FindList(page, limit)
}

func (s *CarService) Update(req *Car) error {
	return s.Update(req)
}
package client

type ClientService struct {
	clientRepo *ClientRepository
}

func NewClientService(clientRepo *ClientRepository) *ClientService {
	return &ClientService{
		clientRepo: clientRepo,
	}
}

func (s *ClientService) FindById(clientId, ownerId string) (*Client, error) {
	return s.clientRepo.FindById(clientId, ownerId)
}

func (s *ClientService) FindListByOwnerId(ownerId string, page, limit int) ([]Client, int64, error) {
	return s.clientRepo.FindByOwnerId(ownerId, page, limit)
}
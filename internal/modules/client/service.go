package client

type ClientService struct {
	clientRepo *ClientRepository
}

func NewClientService(clientRepo *ClientRepository) *ClientService {
	return &ClientService{
		clientRepo: clientRepo,
	}
}

func (s *ClientService) FindById(clientId string) (*Client, error) {
	return s.clientRepo.FindById(clientId)
}
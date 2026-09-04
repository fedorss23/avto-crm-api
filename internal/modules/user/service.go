package user

type UserService struct {
	userRepo *UserRepository
}

func NewUserService(userRepo *UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (r *UserService) Delete(userId string) error {
	return r.userRepo.Delete(userId)
}

func (r *UserService) FindList(page, limit int) ([]User, int64, error) {
	return r.userRepo.FindList(page, limit)
}
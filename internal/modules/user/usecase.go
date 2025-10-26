package user

type UserService = struct {
	repo *UserRepo
}

func NewUserService(repo *UserRepo) *UserService {
	return &UserService{repo: repo}
}

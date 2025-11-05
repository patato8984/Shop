package user

import (
	"errors"

	"github.com/patato8984/Shop/internal/modules/user/model"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *UserRepo
}

func NewUserService(repo *UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterCase(user model.User) error {
	if len(user.Password) < 12 || user.Nickname == "" {
		return errors.New("short password or nickname")
	}
	// add a check if the user is busy
	heshPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	// to add a mailing list for mail, use Kafka
	if err != nil {
		return err
	}
	if err = s.repo.RegisterRepo(user.Mail, user.Name, user.Nickname, string(heshPassword)); err != nil {
		return err
	}
	return nil
}

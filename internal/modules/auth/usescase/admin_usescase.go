package auth_usescase

import (
	"errors"

	"github.com/patato8984/Shop/internal/modules/auth/model"
	repo_user "github.com/patato8984/Shop/internal/modules/auth/repo"
	"golang.org/x/crypto/bcrypt"
)

type AdminServise struct {
	repo *repo_user.AdminRepo
}

func NewAdminServise(repo *repo_user.AdminRepo) *AdminServise {
	return &AdminServise{repo: repo}
}

func (s AdminServise) CreateNewAdmin(user model.User) error {
	if len(user.Password) < 12 || user.Nickname == "" {
		return errors.New("short password or nickname")
	}
	ok, _ := s.repo.SearchNicknameAdmin(user.Nickname)
	if !ok {
		return errors.New("nickname busy")
	}
	heshPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	// to add a mailing list for mail, use Kafka
	if err != nil {
		return err
	}
	if err = s.repo.RegisterAdmin(user.Mail, user.Name, user.Nickname, string(heshPassword)); err != nil {
		return err
	}
	return nil
}

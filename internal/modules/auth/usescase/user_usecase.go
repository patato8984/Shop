package auth_usescase

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/patato8984/Shop/internal/modules/auth/model"
	repo_user "github.com/patato8984/Shop/internal/modules/auth/repo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repo_user.UserRepo
}

func NewUserService(repo *repo_user.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterCase(user model.User) error {
	if len(user.Password) < 12 || user.Nickname == "" {
		return errors.New("short password or nickname")
	}
	ok, _ := s.repo.SearchNickname(user.Nickname)
	if !ok {
		return errors.New("nickname busy")
	}
	heshPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	// to add a mailing list for mail, use Kafka
	if err != nil {
		return err
	}
	if err = s.repo.RegisterUser(user.Mail, user.Name, user.Nickname, string(heshPassword)); err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetToken(nickName, password string) (model.User, error) {
	var user model.User
	users, err := s.repo.GetHashPasswordFromNickname(nickName)
	if err != nil {
		return user, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(users.HeshPassword), []byte(password)); err != nil {
		return user, err
	}
	clearToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{"user_id": users.Id, "role": users.Role, "exp": time.Now().Add(time.Hour * 24).Unix()})
	token, err := clearToken.SignedString([]byte("каки"))
	if err != nil {
		return user, err
	}
	user.CreatedAt = users.CreatedAt
	user.Nickname = nickName
	user.Id = users.Id
	user.Token = token
	return user, nil
}

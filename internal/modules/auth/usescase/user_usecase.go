package auth_usescase

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/patato8984/Shop/internal/modules/auth/model"
	repo_user "github.com/patato8984/Shop/internal/modules/auth/repo"
	shared_events "github.com/patato8984/Shop/internal/shared/events"
	"github.com/patato8984/Shop/internal/shared/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type CartCreatedProvider interface {
	CreatedCart(ctx context.Context, idUser int) (time.Time, error)
}
type UserService struct {
	provider CartCreatedProvider
	repo     *repo_user.UserRepo
	jwtKey   string
	kp       shared_events.EventPublisher
}

func NewUserService(provider CartCreatedProvider, repo *repo_user.UserRepo, jwtKey string, kp shared_events.EventPublisher) *UserService {
	return &UserService{
		provider: provider,
		repo:     repo,
		jwtKey:   jwtKey,
		kp:       kp,
	}
}

func (s *UserService) RegisterCase(ctx context.Context, user model.User) error {
	log := logger.FromContext(ctx)
	if len(user.Password) < 12 || user.Nickname == "" {
		return model.ErrShortPasswordOrNickname
	}
	if err := s.repo.SearchMailUser(ctx, user.Mail); err != nil {
		return err
	}
	if err := s.repo.SearchNicknameUser(ctx, user.Nickname); err != nil {
		return err
	}
	heshPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("error generate hash",
			zap.Error(err),
		)
		return err
	}
	if err = s.repo.RegisterUser(ctx, user.Mail, user.Name, user.Nickname, string(heshPassword)); err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetToken(ctx context.Context, nickName, password string) (model.User, error) {
	var user model.User
	log := logger.FromContext(ctx)
	users, err := s.repo.GetHashPasswordFromNickname(ctx, nickName)
	if err != nil {
		return user, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(users.HeshPassword), []byte(password)); err != nil {
		log.Error("error check hash",
			zap.Error(err),
		)
		return user, model.ErrCheckPassword
	}
	clearToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{"id": users.Id, "role": users.Role, "exp": time.Now().Add(time.Hour * 68).Unix()})
	token, err := clearToken.SignedString([]byte(s.jwtKey))
	if err != nil {
		log.Error("error create jwt",
			zap.Error(err),
		)
		return user, err
	}

	_, err = s.provider.CreatedCart(ctx, users.Id)
	if err != nil {
		return user, err
	}
	user.CreatedAt = users.CreatedAt
	user.Nickname = nickName
	user.Role = users.Role
	user.Id = users.Id
	user.Token = token
	return user, nil
}

package auth_usescase

import (
	"context"
	"log"

	"github.com/patato8984/Shop/internal/modules/auth/model"
	repo_user "github.com/patato8984/Shop/internal/modules/auth/repo"
	"github.com/patato8984/Shop/internal/shared/config"
	"github.com/patato8984/Shop/internal/shared/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AdminService struct {
	repo *repo_user.AdminRepo
}
type SeedService struct {
	admins *config.Admin
	repo   *repo_user.AdminRepo
}

func NewAdminService(repo *repo_user.AdminRepo) *AdminService {
	return &AdminService{repo: repo}
}
func NewSeedService(admins config.Admin, repo *repo_user.AdminRepo) *SeedService {
	return &SeedService{admins: &admins, repo: repo}
}
func (s SeedService) SeedAdmins(ctx context.Context) error {
	for _, admin := range *s.admins {
		if err := s.repo.CheckNicknameBusy(ctx, admin.NickName); err != nil {
			continue
		}
		heshPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
		if err != nil {
			continue
		}
		if err = s.repo.RegisterAdmin(ctx, admin.Mail, admin.Name, admin.NickName, string(heshPassword)); err != nil {
			continue
		}
		log.Printf("add Admin (%s) success", admin.NickName)
	}
	return nil
}
func (s AdminService) CreateNewAdmin(ctx context.Context, user model.User) error {
	log := logger.FromContext(ctx)
	if len(user.Password) < 12 || user.Nickname == "" {
		return model.ErrShortPasswordOrNickname
	}
	if err := s.repo.CheckNicknameBusy(ctx, user.Nickname); err != nil {
		return err
	}
	if err := s.repo.SearchMailUser(ctx, user.Mail); err != nil {
		return err
	}
	heshPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Warn("error generate hash",
			zap.Error(err),
		)
		return err
	}
	if err = s.repo.RegisterAdmin(ctx, user.Mail, user.Name, user.Nickname, string(heshPassword)); err != nil {
		return err
	}
	return nil
}

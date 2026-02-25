package auth_repo

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
	"github.com/patato8984/Shop/internal/modules/auth/model"
	"github.com/patato8984/Shop/internal/shared/logger"
	"go.uber.org/zap"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}
func (r *UserRepo) SearchNicknameUser(ctx context.Context, nickname string) error {
	log := logger.FromContext(ctx)
	var id int
	if err := r.db.QueryRowContext(ctx, "SELECT 1 FROM users WHERE nickname = $1 LIMIT 1", nickname).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		log.Error("error db",
			zap.Error(err),
		)
		return err
	}
	return model.ErrCheckPassword
}
func (r *UserRepo) SearchMailUser(ctx context.Context, mail string) error {
	log := logger.FromContext(ctx)
	var id int
	if err := r.db.QueryRowContext(ctx, "SELECT 1 FROM users WHERE mail = $1 LIMIT 1", mail).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		log.Error("error db",
			zap.Error(err),
		)
		return err
	}
	return model.ErrMailBusy
}
func (r *UserRepo) RegisterUser(ctx context.Context, mail, name, nickname, hashPassword string) error {
	log := logger.FromContext(ctx)
	_, err := r.db.ExecContext(ctx, "INSERT INTO users (mail, name, nickName, password, role) VALUES ($1, $2, $3, $4, $5)", mail, name, nickname, hashPassword, "user")
	if err != nil {
		log.Error("error db",
			zap.Error(err),
		)
		return err
	}
	return nil
}
func (r *UserRepo) GetHashPasswordFromNickname(ctx context.Context, nickname string) (model.ResponseAuthentication, error) {
	var user model.ResponseAuthentication
	log := logger.FromContext(ctx)
	if err := r.db.QueryRowContext(ctx, "SELECT id, password, role, created_at FROM users WHERE nickName = $1", nickname).Scan(&user.Id, &user.HeshPassword, &user.Role, &user.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, model.ErrUserNotFound
		}
		log.Error("error db",
			zap.Error(err),
		)
		return user, err
	}
	return user, nil
}

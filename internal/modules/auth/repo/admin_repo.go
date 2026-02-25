package auth_repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/patato8984/Shop/internal/modules/auth/model"
	"github.com/patato8984/Shop/internal/shared/logger"
	"go.uber.org/zap"
)

type AdminRepo struct {
	db *sql.DB
}

func NewAdminRepo(db *sql.DB) *AdminRepo {
	return &AdminRepo{db: db}
}
func (r *AdminRepo) CheckNicknameBusy(ctx context.Context, nickname string) error {
	var id int
	log := logger.FromContext(ctx)
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
func (r *AdminRepo) SearchMailUser(ctx context.Context, mail string) error {
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
func (r *AdminRepo) RegisterAdmin(ctx context.Context, mail, name, nickname, hashPassword string) error {
	log := logger.FromContext(ctx)
	if _, err := r.db.ExecContext(ctx, "INSERT INTO users (name, nickName, mail, password, role) VALUES ($1, $2, $3, $4, $5) ", name, nickname, mail, hashPassword, "admin"); err != nil {
		log.Error("error db",
			zap.Error(err),
		)
		return err
	}
	return nil
}

package auth_repo

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/patato8984/Shop/internal/modules/auth/model"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}
func (r *UserRepo) SearchNickname(nickname string) (bool, error) {
	var ok bool
	err := r.db.QueryRow("SELECT 1 FROM users WHERE nickName = $1", nickname).Scan(&ok)
	if err != nil {
		return false, err
	}
	return ok, nil
}
func (r *UserRepo) RegisterUser(mail, name, nickname, hashPassword string) error {
	_, err := r.db.Exec("INSERT INTO users (user_mail,user_name, user_nickName, user_password, role) VALUES ($1, $2, $3, $4, %5)", mail, name, nickname, hashPassword, "user")
	if err != nil {
		return err
	}
	return nil
}
func (r *UserRepo) GetHashPasswordFromNickname(nickname string) (model.ResponseAuthentication, error) {
	var user model.ResponseAuthentication
	if err := r.db.QueryRow("SELECT id, password, role, created_at FROM users WHERE nickName = $1", nickname).Scan(&user.Id, &user.HeshPassword, &user.Role, &user.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errors.New("user not found")
		}
		return user, fmt.Errorf("db error: %w", err)
	}
	return user, nil
}

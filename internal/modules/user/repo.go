package user

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/patato8984/Shop/internal/modules/user/model"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}
func (r *UserRepo) SearchNickname(nickname string) (bool, error) {
	var ok bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM user WHERE nickName = $1)", nickname).Scan(&ok)
	if err != nil {
		return false, err
	}
	return ok, nil
}
func (r *UserRepo) RegisterRepo(mail, name, nickname, hashPassword string) error {
	_, err := r.db.Exec("INSERT INTO user (user_mail,user_name, user_nickName, user_password) VALUES ($1, $2, $3, $4)", mail, name, nickname, hashPassword)
	if err != nil {
		return err
	}
	return nil
}
func (r *UserRepo) GetHashPasswordFromNickname(nickname string) (model.HashPasswordAndId, error) {
	var user model.HashPasswordAndId
	if err := r.db.QueryRow("SELECT id, password FROM user WHERE nickName = $1", nickname).Scan(&user.Id, &user.HeshPassword); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errors.New("user not found")
		}
		return user, fmt.Errorf("db error: %w", err)
	}
	return user, nil
}

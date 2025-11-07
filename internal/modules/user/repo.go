package user

import (
	"database/sql"

	_ "github.com/lib/pq"
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

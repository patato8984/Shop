package user

import (
	"database/sql"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}
func (r *UserRepo) RegisterRepo(mail, name, nickname, hashPassword string) error {
	_, err := r.db.Exec("INSERT INTO user (user_mail,user_name, user_nickName, user_password) VALUES ($1, $2, $3, $4)", mail, name, nickname, hashPassword)
	if err != nil {
		return err
	}
	return nil
}

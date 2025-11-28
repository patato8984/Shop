package auth_repo

import (
	"database/sql"
)

type AdminRepo struct {
	db *sql.DB
}

func NewAdminRepo(db *sql.DB) *AdminRepo {
	return &AdminRepo{db: db}
}
func (r *AdminRepo) SearchNicknameAdmin(nickname string) (bool, error) {
	var ok bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE nickName = $1)", nickname).Scan(&ok)
	if err != nil {
		return false, err
	}
	return ok, nil
}
func (r *AdminRepo) RegisterAdmin(mail, name, nickname, hashPassword string) error {
	if _, err := r.db.Exec("INSERT INTO users (name, nickName, mail, password, role)", name, nickname, mail, hashPassword, "admin"); err != nil {
		return err
	}
	return nil
}

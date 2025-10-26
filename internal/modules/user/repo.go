package user

import (
	"database/sql"
)

type UserRepo = struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

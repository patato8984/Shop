package order_repo

import "database/sql"

type OrderUserRepo struct {
	db *sql.DB
}

func NewOrderUserRepo(db *sql.DB) *OrderUserRepo {
	return &OrderUserRepo{db: db}
}

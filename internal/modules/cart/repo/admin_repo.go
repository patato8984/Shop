package cart_repo

import "database/sql"

type CartAdminRepo struct {
	db *sql.DB
}

func NewCatalogAmdinRepo(db *sql.DB) *CartAdminRepo {
	return &CartAdminRepo{db: db}
}

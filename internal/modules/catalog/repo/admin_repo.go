package repo

import (
	"database/sql"
)

type CatalogAdminRepo struct {
	db *sql.DB
}

func NewCatalogAdminRepo(db *sql.DB) *CatalogAdminRepo {
	return &CatalogAdminRepo{db: db}
}
func (r CatalogAdminRepo) CreateOrSearchProduct(product string) (int, error) {
	var id int
	if err := r.db.QueryRow("INSERT INTO products (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id", product).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

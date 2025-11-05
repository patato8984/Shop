package repo

import (
	"database/sql"

	modelscatalog "github.com/patato8984/Shop/internal/modules/catalog/model"
)

type CatalogRepo struct {
	db *sql.DB
}

func NewCatalogRepo(db *sql.DB) *CatalogRepo {
	return &CatalogRepo{db: db}
}
func (r CatalogRepo) GetAll() ([]modelscatalog.Product, error) {
	rows, err := r.db.Query("SELECT id, name, stock, price FROM products WHERE stock > 0")
	if err != nil {
		return []modelscatalog.Product{}, err
	}
	var products []modelscatalog.Product
	defer rows.Close()
	for rows.Next() {
		var product modelscatalog.Product
		if err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Price); err != nil {
			return products, err
		}
		products = append(products, product)
	}
	return products, err
}

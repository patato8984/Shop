package repo

import (
	"database/sql"

	"github.com/patato8984/Shop/internal/modules/catalog/model"
)

type CatalogRepo struct {
	db *sql.DB
}

func NewCatalogRepo(db *sql.DB) *CatalogRepo {
	return &CatalogRepo{db: db}
}
func (r CatalogRepo) GetAll() ([]model.Product, error) {
	rows, err := r.db.Query("SELECT id, name FROM products")
	if err != nil {
		return []model.Product{}, err
	}
	var products []model.Product
	defer rows.Close()
	for rows.Next() {
		var product model.Product
		if err := rows.Scan(&product.Id, &product.Name); err != nil {
			return products, err
		}
		products = append(products, product)
	}
	return products, err
}
func (r CatalogRepo) GetProduct(id int) (model.Product, error) {
	rows, err := r.db.Query("SELECT p.id, p.name, s.id AS id, s.products_id, s.storage, s.colour, s.price, s.stock FROM products p JOIN skus s ON p.id = s.products_id WHERE p.id = $1", id)
	var product model.Product
	if err != nil {
		return product, err
	}
	defer rows.Close()
	product.SKUs = []model.SKU{}
	for rows.Next() {
		var sku model.SKU
		var productID int
		var productName string
		if err := rows.Scan(&productID, &productName, &sku.Id, &sku.Storage, &sku.Colour, &sku.Price, &sku.Stock); err != nil {
			return product, err
		}
		product.Id = productID
		product.Name = productName
		product.SKUs = append(product.SKUs, sku)
	}
	return product, nil
}

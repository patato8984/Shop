package catalog_repo

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/patato8984/Shop/internal/modules/catalog/model"
)

type CatalogRepo struct {
	db *sql.DB
}

func NewCatalogRepo(db *sql.DB) *CatalogRepo {
	return &CatalogRepo{db: db}
}
func (r CatalogRepo) GetAll() ([]model.Product, error) {
	rows, err := r.db.Query("SELECT id, name FROM products WHERE delete_at IS nil")
	if err != nil {
		return []model.Product{}, fmt.Errorf("db error: %w", err)
	}
	var products []model.Product
	defer rows.Close()
	for rows.Next() {
		var product model.Product
		if err := rows.Scan(&product.Id, &product.Name); err != nil {
			return products, fmt.Errorf("db error: %w", err)
		}
		products = append(products, product)
	}
	return products, nil
}
func (r CatalogRepo) GetSkus(id int) (model.SKU, error) {
	var skus model.SKU
	err := r.db.QueryRow("SELECT id, products_id, storage, colour, price, stock FROM skus WHERE id = $id AND deleted_at IS NILL", id).Scan(&skus.Id, &skus.Product_id, &skus.Storage, &skus.Colour, &skus.Price, &skus.Stock)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return skus, errors.New("skus not found")
		}
		return skus, fmt.Errorf("db error: %w", err)
	}
	return skus, nil
}
func (r CatalogRepo) GetAllSkus(id int) (model.Product, error) {
	rows, err := r.db.Query("SELECT p.id, p.name, s.id AS id, s.products_id, s.storage, s.colour, s.price, s.stock FROM products p JOIN skus s ON p.id = s.products_id WHERE p.id = $1", id)
	var product model.Product
	if err != nil {
		return product, fmt.Errorf("db error: %w", err)
	}
	defer rows.Close()
	product.SKUs = []model.SKU{}
	for rows.Next() {
		var sku model.SKU
		var productID int
		var productName string
		if err := rows.Scan(&productID, &productName, &sku.Id, &sku.Storage, &sku.Colour, &sku.Price, &sku.Stock); err != nil {
			return product, fmt.Errorf("db error: %w", err)
		}
		product.Id = productID
		product.Name = productName
		product.SKUs = append(product.SKUs, sku)
	}
	return product, nil
}
func (r CatalogRepo) GetPrice(idSkus int) (float64, error) {
	var price float64
	err := r.db.QueryRow("SELECT price FROM skus WHERE id = $1", idSkus).Scan(&price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("skus not faund")
		}
		return 0, fmt.Errorf("error db: %s", err)
	}
	return price, nil
}
func (r CatalogRepo) GetStock(idSkus int) (int, error) {
	var stock int
	err := r.db.QueryRow("SELECT stock FROM skus WHERE id = $1", idSkus).Scan(&stock)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("skus not faund")
		}
		return 0, fmt.Errorf("error db: %s", err)
	}
	return stock, nil
}

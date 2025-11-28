package catalog_repo

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/patato8984/Shop/internal/modules/catalog/model"
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
		return 0, fmt.Errorf("db error: %w", err)
	}
	return id, nil
}
func (r CatalogAdminRepo) DelProduct(id int) (int, error) {
	var deletedID int
	err := r.db.QueryRow("UPDATE products SET deleted_at = NEW() WHERE id = $1 RETURNING id", id).Scan(&deletedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("product not found")
		}
		return 0, fmt.Errorf("db error: %w", err)
	}
	return deletedID, nil
}

func (r CatalogAdminRepo) CreateSkus(id int, sku model.SKU) (model.Product, error) {
	var product model.Product
	var createdSku model.SKU
	var sku_id int
	err := r.db.QueryRow("INSERT INTO skus (products_id, storage, colour, price) VALUES ($1, $2, $3, $4) RETURNING id", id, sku.Storage, sku.Colour, sku.Price).Scan(&sku_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return product, errors.New("product not found")
		}
		return product, fmt.Errorf("db error: %w", err)
	}
	if err := r.db.QueryRow("SELECT p.id, p.name, s.id AS sku_id, s.products_id, s.storage, s.colour, s.price, s.stock FROM products p JOIN skus s ON s.products_id = p.id WHERE p.id = $1 AND s.id = $2", id, sku_id).Scan(
		&product.Id, &product.Name, &createdSku.Id,
		&createdSku.Product_id, &createdSku.Storage, &createdSku.Colour,
		&createdSku.Price, &createdSku.Stock); err != nil {
		return product, err
	}
	product.SKUs = append(product.SKUs, createdSku)
	return product, nil
}
func (r CatalogAdminRepo) AddStock(stock, id int) (model.SKU, error) {
	var sku model.SKU
	err := r.db.QueryRow("UPDATE skus SET stock = stock + $1, appdate_at = NEW() WHERE products_id = $2 RETURNING products_id, price, stock", stock, id).Scan(&sku.Product_id, &sku.Price, &sku.Stock)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sku, errors.New("skus not found")
		}
		return sku, fmt.Errorf("db error: %w", err)
	}
	sku.Id = id
	return sku, nil
}
func (r CatalogAdminRepo) DelSkus(id int) (int, error) {
	var deletedID int
	err := r.db.QueryRow("UPDATE skus SET deleted_at = NEW() WHERE id = $1 RETURNING id", id).Scan(&deletedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("skus not found")
		}
		return 0, fmt.Errorf("db error: %w", err)
	}
	return deletedID, nil
}

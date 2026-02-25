package catalog_repo

import (
	"context"
	"database/sql"
	"errors"

	catalog_model "github.com/patato8984/Shop/internal/modules/catalog/model"
	"github.com/patato8984/Shop/internal/shared/dto"
	"github.com/patato8984/Shop/internal/shared/logger"
	"go.uber.org/zap"
)

type CatalogRepo struct {
	db *sql.DB
}

func NewCatalogRepo(db *sql.DB) *CatalogRepo {
	return &CatalogRepo{db: db}
}
func (r *CatalogRepo) GetAll(ctx context.Context) (*[]catalog_model.Product, error) {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	rows, err := db.QueryContext(ctx, "SELECT id, name FROM products WHERE deleted_at IS NULL")
	if err != nil {
		log.Error("error db",
			zap.Error(err),
		)
		return &[]catalog_model.Product{}, err
	}
	var products []catalog_model.Product
	defer rows.Close()
	for rows.Next() {
		var product catalog_model.Product
		if err := rows.Scan(&product.Id, &product.Name); err != nil {
			log.Error("error rows db",
				zap.Error(err),
			)
			return &products, err
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		log.Error("error rows db",
			zap.Error(err),
		)
		return nil, err
	}
	return &products, nil
}
func (r *CatalogRepo) GetSkus(ctx context.Context, id int) (*catalog_model.SKU, error) {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	var skus catalog_model.SKU
	err := db.QueryRowContext(ctx, "SELECT id, products_id, storage, colour, price, stock FROM skus WHERE id = $1 AND deleted_at IS NULL", id).Scan(&skus.Id, &skus.Product_id, &skus.Storage, &skus.Colour, &skus.Price, &skus.Stock)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &skus, catalog_model.ErrSkuNotFound
		}
		log.Error("error db",
			zap.Error(err),
		)
		return &skus, err
	}
	return &skus, nil
}
func (r *CatalogRepo) GetAllSkus(ctx context.Context, id int) (*catalog_model.Product, error) {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	rows, err := db.QueryContext(ctx, "SELECT p.id, p.name, s.id AS id, s.storage, s.colour, s.price, s.stock FROM products p JOIN skus s ON p.id = s.products_id WHERE p.id = $1", id)
	var product catalog_model.Product
	if err != nil {
		log.Error("error db",
			zap.Error(err),
		)
		return &product, err
	}
	defer rows.Close()
	product.SKUs = []catalog_model.SKU{}
	for rows.Next() {
		var sku catalog_model.SKU
		var productID int
		var productName string
		if err := rows.Scan(&productID, &productName, &sku.Id, &sku.Storage, &sku.Colour, &sku.Price, &sku.Stock); err != nil {
			log.Error("error rows db",
				zap.Error(err),
			)
			return &product, err
		}
		product.Id = productID
		product.Name = productName
		product.SKUs = append(product.SKUs, sku)
	}
	if err := rows.Err(); err != nil {
		log.Error("error rows db",
			zap.Error(err),
		)
		return &product, err
	}
	return &product, nil
}
func (r *CatalogRepo) SearchIdProduct(ctx context.Context, idSkus int) error {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	var id int
	err := db.QueryRowContext(ctx, "SELECT id FROM skus WHERE id = $1", idSkus).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return catalog_model.ErrSkuNotFound
		}
		log.Error("error db",
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (r *CatalogRepo) GetStock(ctx context.Context, idSkus int) (int, error) {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	var stock int
	err := db.QueryRowContext(ctx, "SELECT stock FROM skus WHERE id = $1", idSkus).Scan(&stock)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, catalog_model.ErrSkuNotFound
		}
		log.Error("error db",
			zap.Error(err),
		)
		return 0, err
	}
	return stock, nil
}
func (r *CatalogRepo) WorkerGetPrice(ctx context.Context, idSkus int) (float64, error) {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	var price float64
	err := db.QueryRowContext(ctx, "SELECT price FROM skus WHERE id = $1", idSkus).Scan(&price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return price, catalog_model.ErrSkuNotFound
		}
		log.Error("error db",
			zap.Error(err),
		)
		return price, err
	}
	return price, nil

}

package catalog_repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	catalog_model "github.com/patato8984/Shop/internal/modules/catalog/model"
	"github.com/patato8984/Shop/internal/shared/dto"
	"github.com/patato8984/Shop/internal/shared/logger"
	"github.com/patato8984/Shop/pkg/ctxkey"
	"go.uber.org/zap"
)

type CatalogAdminRepo struct {
	db *sql.DB
}
type TxManager struct {
	db *sql.DB
}

func NewCatalogTxManager(db *sql.DB) *TxManager {
	return &TxManager{db: db}
}
func NewCatalogAdminRepo(db *sql.DB) *CatalogAdminRepo {
	return &CatalogAdminRepo{db: db}
}
func (r *TxManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) (any, error)) (any, error) {
	log := logger.FromContext(ctx)
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error("error start transaction db",
			zap.Error(err),
		)
		return nil, err
	}
	ctxWitchTx := context.WithValue(ctx, ctxkey.TransactionKey, tx)
	result, err := fn(ctxWitchTx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Error("error rollback db",
				zap.Error(err),
			)
			return nil, err
		}
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		log.Error("error create commit db",
			zap.Error(err),
		)
		return nil, err
	}
	return result, nil
}

func (r CatalogAdminRepo) CreateOrSearchProduct(ctx context.Context, product string) (int, error) {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	var id int
	if err := db.QueryRowContext(ctx, "INSERT INTO products (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id", product).Scan(&id); err != nil {
		log.Error("error db",
			zap.Error(err),
		)
		return 0, err
	}
	return id, nil
}
func (r CatalogAdminRepo) DelProduct(ctx context.Context, id int) (int, error) {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	var deletedID int
	err := db.QueryRowContext(ctx, "UPDATE products SET deleted_at = NOW() WHERE id = $1 RETURNING id", id).Scan(&deletedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, catalog_model.ErrProductNotFound
		}
		log.Error("error db",
			zap.Error(err),
		)
		return 0, err
	}
	return deletedID, nil
}

func (r CatalogAdminRepo) CreateSkus(ctx context.Context, id int, sku catalog_model.SKU) (catalog_model.Product, error) {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	var product catalog_model.Product
	var createdSku catalog_model.SKU
	var sku_id int
	err := db.QueryRowContext(ctx, "INSERT INTO skus (products_id, storage, colour, price, stock) VALUES ($1, $2, $3, $4, $5) RETURNING id", id, sku.Storage, sku.Colour, sku.Price, sku.Stock).Scan(&sku_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return product, catalog_model.ErrProductNotFound
		}
		log.Error("error db",
			zap.Error(err),
		)
		return product, err
	}
	if err := db.QueryRowContext(ctx, "SELECT p.id, p.name, s.id AS sku_id, s.products_id, s.storage, s.colour, s.price, s.stock FROM products p JOIN skus s ON s.products_id = p.id WHERE p.id = $1 AND s.id = $2", id, sku_id).Scan(
		&product.Id, &product.Name, &createdSku.Id,
		&createdSku.Product_id, &createdSku.Storage, &createdSku.Colour,
		&createdSku.Price, &createdSku.Stock); err != nil {
		log.Error("error db",
			zap.Error(err),
		)
		return product, err
	}
	product.SKUs = append(product.SKUs, createdSku)
	return product, nil
}
func (r CatalogAdminRepo) GetStock(ctx context.Context, id int) (int, error) {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	var stock int
	err := db.QueryRowContext(ctx, "SELECT stock FROM id = $1 LIMIT 1", id).Scan(&stock)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, catalog_model.ErrProductNotFound
		}
		log.Error("error db",
			zap.Error(err),
		)
		return 0, err
	}
	return stock, nil
}
func (r CatalogAdminRepo) UpdateStock(ctx context.Context, delta, id int) (int, int, time.Time, error) {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	var newStock int
	var update_at time.Time
	var idProduct int
	err := db.QueryRowContext(ctx, "UPDATE skus SET stock = stock + $1, update_at = NOW() WHERE id = $2 RETURNING stock, update_at, products_id", delta, id).Scan(&newStock, &update_at, &idProduct)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return idProduct, newStock, update_at, catalog_model.ErrSkuNotFound
		}
		log.Error("error db",
			zap.Error(err),
		)
		return idProduct, newStock, update_at, err
	}
	return idProduct, newStock, update_at, nil
}
func (r CatalogAdminRepo) UpdatePriceToSkus(ctx context.Context, idSkus int, newPrice float64) (int, float64, error) {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	var price float64
	var idProduct int
	err := db.QueryRowContext(ctx, "UPDATE skus SET price = $1 WHERE id = $2 RETURNING id, products_id", newPrice, idSkus).Scan(&price, &idProduct)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return idProduct, price, catalog_model.ErrSkuNotFound
		}
		log.Error("error db",
			zap.Error(err),
		)
		return idProduct, price, err
	}
	return idProduct, price, nil
}
func (r CatalogAdminRepo) DelSkus(ctx context.Context, id int) (int, int, error) {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	var deletedID int
	var idProduct int
	err := db.QueryRowContext(ctx, "UPDATE skus SET deleted_at = NOW() WHERE id = $1 RETURNING id, products_id", id).Scan(&deletedID, &idProduct)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, catalog_model.ErrSkuNotFound
		}
		log.Error("error db",
			zap.Error(err),
		)
		return 0, 0, err
	}
	return idProduct, deletedID, nil
}

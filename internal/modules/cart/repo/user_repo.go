package cart_repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	cart_model "github.com/patato8984/Shop/internal/modules/cart/model"
	"github.com/patato8984/Shop/internal/shared/dto"
	"github.com/patato8984/Shop/internal/shared/logger"
	"github.com/patato8984/Shop/pkg/ctxkey"
	"go.uber.org/zap"
)

type CartRepo struct {
	db *sql.DB
}
type TxManager struct {
	db *sql.DB
}

func NewCartRepo(db *sql.DB) *CartRepo {
	return &CartRepo{db: db}
}
func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{db: db}
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
func (r *CartRepo) SearchCart(ctx context.Context) (int, error) {
	log := logger.FromContext(ctx)
	idUser := ctx.Value(ctxkey.UserIDKey)
	db := dto.Getter(ctx, r.db)
	var id int
	err := db.QueryRowContext(ctx, "SELECT id FROM cart WHERE id_user = $1 AND status = $2 LIMIT 1", idUser, "active").Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, cart_model.ErrCartNotFound
		}
		log.Error("error db",
			zap.Error(err),
		)
		return 0, err
	}
	return id, nil
}
func (r *CartRepo) GetNumberItemsSkus(ctx context.Context, idCart int) (int, error) {
	log := logger.FromContext(ctx)
	var quantity int
	db := dto.Getter(ctx, r.db)
	err := db.QueryRowContext(ctx, "SELECT COUNT(id_skus) FROM cart_items WHERE id_cart = $1", idCart).Scan(&quantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, cart_model.ErrCartNotFound
		}
		log.Error("error db",
			zap.Error(err),
		)
		return 0, err
	}
	return quantity, nil
}
func (r *CartRepo) CreatedCart(ctx context.Context, idUser int) (time.Time, error) {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	var update_at time.Time
	err := db.QueryRowContext(ctx, "INSERT INTO cart (id_user, created_at) VALUES ($1, NOW()) RETURNING update_at", idUser)
	if err.Err() != nil {
		log.Error("error db",
			zap.Error(err.Err()),
		)
		return update_at, err.Err()
	}
	return update_at, nil
}
func (r *CartRepo) GetIdCartItems(ctx context.Context, idCart, idSkus int) (int, error) {
	log := logger.FromContext(ctx)
	var idCartItems int
	db := dto.Getter(ctx, r.db)
	if err := db.QueryRowContext(ctx, "SELECT id FROM cart_items WHERE id_cart = $1 AND id_skus = $2 LIMIT 1 RETURNING id", idCart, idSkus).Scan(&idCartItems); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, cart_model.ErrCartItemsNotFund
		}
		log.Error("error db",
			zap.Error(err),
		)
		return 0, err
	}
	return idCartItems, nil
}
func (r *CartRepo) GetAllIdItemsCart(ctx context.Context, idCart int) ([]int, error) {
	log := logger.FromContext(ctx)
	var allIdItems []int
	db := dto.Getter(ctx, r.db)
	rows, err := db.QueryContext(ctx, "SELECT id FROM cart_items WHERE id_cart = $1", idCart)
	if err != nil {
		log.Error("error db",
			zap.Error(err),
		)
		return allIdItems, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Error("error scan db",
				zap.Error(err),
			)
			return allIdItems, err
		}
		allIdItems = append(allIdItems, id)
	}
	return allIdItems, nil
}
func (r *CartRepo) GetAllCostProduct(ctx context.Context, idCart int) (*float64, error) {
	log := logger.FromContext(ctx)
	var price float64
	db := dto.Getter(ctx, r.db)
	err := db.QueryRowContext(ctx, "SELECT COALESCE(SUM(price_snapshot), 0) FROM cart_items WHERE id_cart = $1", idCart).Scan(&price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &price, cart_model.ErrCartItemsNotFund
		}
		log.Error("error db",
			zap.Error(err),
		)
		return &price, err
	}
	return &price, nil
}
func (r *CartRepo) DelProduct(ctx context.Context, idCart, idSkus int) (int, error) {
	log := logger.FromContext(ctx)
	var id int
	db := dto.Getter(ctx, r.db)
	if err := db.QueryRowContext(ctx, "DELETE FROM cart_items WHERE id_cart = $1 AND id_skus = $2 RETURNING id").Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return id, cart_model.ErrCartItemsNotFund
		}
		log.Error("error db",
			zap.Error(err),
		)
		return id, err
	}
	return id, nil
}
func (r *CartRepo) ClearAllItems(ctx context.Context, idCart int) error {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	_, err := db.ExecContext(ctx, "DELETE from cart_items WHERE id_cart = $1", idCart)
	if err != nil {
		log.Error("error db",
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (r *CartRepo) GetCart(ctx context.Context, idCart int) (*cart_model.Cart, error) {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	rows, err := db.QueryContext(ctx, "SELECT c.id AS cart_id, c.id_user, c.created_at, c.update_at, c.status, i.id AS id_cartItems, i.id_cart, i.id_skus, i.quantity, i.price_snapshot FROM cart c JOIN cart_items i ON c.id = i.id_cart WHERE c.id = $1", idCart)
	var cart cart_model.Cart
	if err != nil {
		log.Error("error db",
			zap.Error(err),
		)
		return &cart, err
	}
	defer rows.Close()
	cart.Cart_items = []cart_model.Items{}
	for rows.Next() {
		items := cart_model.Items{}
		var idCart int
		var idUser int
		var created_at time.Time
		var update_at sql.NullTime
		var status string
		if err := rows.Scan(&idCart, &idUser, &created_at, &update_at, &status, &items.Id, &items.Id_cart, &items.Id_skus, &items.Quantity, &items.Price_snapshot); err != nil {
			log.Error("error rows db",
				zap.Error(err),
			)
			return &cart, err
		}
		cart.PriceAllProductCart += items.Price_snapshot
		cart.Id_user = idUser
		cart.Created_at = created_at
		if update_at.Valid {
			cart.Update_at = update_at
		}
		cart.Status = status
		cart.Cart_items = append(cart.Cart_items, items)
	}
	if err := rows.Err(); err != nil {
		log.Error("error db",
			zap.Error(err),
		)
		return &cart, err
	}
	return &cart, nil
}
func (r *CartRepo) AddUpdate_at(ctx context.Context, idCart int) (time.Time, error) {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	var time time.Time
	err := db.QueryRowContext(ctx, "UPDATE cart SET update_at = NOW() WHERE id = $1 RETURNING update_at", idCart).Scan(&time)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return time, cart_model.ErrCartNotFound
		}
		log.Error("error db",
			zap.Error(err),
		)
		return time, err
	}
	return time, nil
}

func (r CartRepo) AddProduct(ctx context.Context, idCart, idSkus, quantity int, price_snapshot float64) (int, error) {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	var idCartItems int
	err := db.QueryRowContext(ctx, "INSERT INTO cart_items (id_cart, id_skus, quantity, price_snapshot) VALUES ($1, $2, $3, $4) RETURNING id", idCart, idSkus, quantity, price_snapshot).Scan(&idCartItems)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, cart_model.ErrCartNotFound
		}
		log.Error("error db",
			zap.Error(err),
		)
		return 0, err
	}
	return idCartItems, nil
}
func (r *CartRepo) GetAllIdSkusAndQuantity(ctx context.Context, idCart int) (map[int]int, error) {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	var idAndQuantity = make(map[int]int)
	rows, err := db.QueryContext(ctx, "SELECT id_skus, quantity FROM cart_items WHERE id_cart = $1", idCart)
	if err != nil {
		log.Error("error db",
			zap.Error(err),
		)
		return idAndQuantity, err
	}
	defer rows.Close()
	for rows.Next() {
		var idSkus int
		var quantity int
		err := rows.Scan(&idSkus, &quantity)
		if err != nil {
			log.Error("error rows db",
				zap.Error(err),
			)
			return idAndQuantity, err
		}
		idAndQuantity[idSkus] = quantity
	}
	if err := rows.Err(); err != nil {
		log.Error("error rows db",
			zap.Error(err),
		)
		return idAndQuantity, err
	}
	return idAndQuantity, nil
}
func (r *CartRepo) UpdateQuantity(ctx context.Context, idCart, idSkus, delta int) error {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	_, err := db.ExecContext(ctx, "UPDATE cart_items SET quantity = quantity + $1 WHERE id_cart = $2 AND id_skus = $3", delta, idCart, delta)
	if err != nil {
		log.Error("error db",
			zap.Error(err),
		)
		return err
	}
	return nil
}
func (r *CartRepo) GetQuantity(ctx context.Context, idCartItems, idSkus int) (int, error) {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	var quantity int
	err := db.QueryRowContext(ctx, "SELECT quantity FROM cart_items WHERE id = $1 AND id_skus = $2 RETURNING quantity", idCartItems, idSkus).Scan(&quantity)
	if err != nil {
		log.Error("error db",
			zap.Error(err),
		)
		return 0, err
	}
	return quantity, nil
}
func (r *CartRepo) WorkerGetAllProductId(ctx context.Context) (map[int][]int, error) {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	var id_Product = make(map[int][]int)
	rows, err := db.QueryContext(ctx, "SELECT c.id, i.id_skus FROM cart c JOIN cart_items i ON c.id = i.id_cart WHERE c.status = 'active'")
	if err != nil {
		log.Error("error db",
			zap.Error(err),
		)
		return id_Product, err
	}
	defer rows.Close()
	for rows.Next() {
		var id_cart int
		var id_skus int
		if err := rows.Scan(&id_cart, &id_skus); err != nil {
			log.Error("error rows db",
				zap.Error(err),
			)
			return id_Product, err
		}
		id_Product[id_skus] = append(id_Product[id_skus], id_cart)
	}
	if err := rows.Err(); err != nil {
		log.Error("error rows db",
			zap.Error(err),
		)
		return id_Product, err
	}
	return id_Product, nil
}
func (r *CartRepo) WorkerUpdatePrice(ctx context.Context, id_cart, id_skus int, price float64) error {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	var id int
	err := db.QueryRowContext(ctx, "UPDATE cart_items SET price_snapshot = $1::numeric * quantity WHERE id_cart = $2 AND id_skus = $3 RETURNING id", price, id_cart, id_skus).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return cart_model.ErrCartItemsNotFund
		}
		log.Error("error db",
			zap.Error(err),
		)
		return err
	}
	return nil
}

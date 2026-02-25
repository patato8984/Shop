package order_repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	order_model "github.com/patato8984/Shop/internal/modules/order/model"
	"github.com/patato8984/Shop/internal/shared/dto"
	"github.com/patato8984/Shop/internal/shared/logger"
	"github.com/patato8984/Shop/pkg/ctxkey"
	"go.uber.org/zap"
)

type OrderUserRepo struct {
	db *sql.DB
}
type TxManager struct {
	db *sql.DB
}

func NewOrderUserRepo(db *sql.DB) *OrderUserRepo {
	return &OrderUserRepo{db: db}
}
func NewOrderTxManager(db *sql.DB) *TxManager {
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
			return nil, errors.New("error Rollback")
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

func (r *OrderUserRepo) SearchAllOrders(ctx context.Context) ([]int, error) {
	log := logger.FromContext(ctx)
	idUser := ctx.Value(ctxkey.UserIDKey)
	db := dto.Getter(ctx, r.db)
	var idOrders []int
	rows, err := db.QueryContext(ctx, "SELECT id FROM orders WHERE id_user = $1 RETURNING id", idUser)
	if err != nil {
		log.Error("error db",
			zap.Error(err),
		)
		return idOrders, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Error("error rows db",
				zap.Error(err),
			)
			return idOrders, err
		}
		idOrders = append(idOrders, id)
	}
	if err := rows.Err(); err != nil {
		log.Error("error rows db",
			zap.Error(err),
		)
		return idOrders, err
	}
	return idOrders, nil
}
func (r *OrderUserRepo) GetAllOrdersBaseInfo(ctx context.Context) ([]order_model.Order, error) {
	log := logger.FromContext(ctx)
	var orders []order_model.Order
	db := dto.Getter(ctx, r.db)
	idUser := ctx.Value(ctxkey.UserIDKey)
	rows, err := db.QueryContext(ctx, "SELECT id, addres, statuse, total_amount, price_all_products, created_at FROM orders WHERE id_user = $1", idUser)
	if err != nil {
		log.Error("error db",
			zap.Error(err),
		)
		return orders, err
	}
	defer rows.Close()
	for rows.Next() {
		var order order_model.Order
		if err := rows.Scan(&order.Id, &order.Address, &order.Status, &order.Total_amount, &order.Price_all_products, &order.Created_at); err != nil {
			log.Error("error rows db",
				zap.Error(err),
			)
			return orders, err
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		log.Error("error rows db",
			zap.Error(err),
		)
		return orders, err
	}
	return orders, nil
}
func (r *OrderUserRepo) AddOrder(ctx context.Context, address, message string) (int, time.Time, error) {
	log := logger.FromContext(ctx)
	idUser := ctx.Value(ctxkey.UserIDKey)
	db := dto.Getter(ctx, r.db)
	var timeCreated time.Time
	var idOrder int
	err := db.QueryRowContext(ctx, "INSERT INTO orders (id_user, addres, messages) VALUES ($1, $2, $3) RETURNING created_at, id", idUser, address, message).Scan(&timeCreated, &idOrder)
	if err != nil {
		log.Error("error db",
			zap.Error(err),
		)
		return 0, timeCreated, err
	}
	return idOrder, timeCreated, nil
}

func (r *OrderUserRepo) AddBankPaymentId(ctx context.Context, idOrder int, idPayment string) error {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	var newIdOrder int
	err := db.QueryRowContext(ctx, "UPDATE orders SET bank_payment_id = $1 WHERE id = $2 RETURNING id", idPayment, idOrder).Scan(&newIdOrder)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return order_model.ErrOrderNotFound
		}
		log.Error("error db",
			zap.Error(err),
		)
		return err
	}
	return nil
}
func (r *OrderUserRepo) GetTotalQuantity(ctx context.Context, idOrder int) (int, error) {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	var quantityAllProduct int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM orders_items WHERE id_order = $1", idOrder).Scan(&quantityAllProduct)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, order_model.ErrOrderNotFound
		}
		log.Error("error db",
			zap.Error(err),
		)
		return 0, err
	}
	return quantityAllProduct, nil
}
func (r *OrderUserRepo) GetOrder(ctx context.Context, idOrder int) (order_model.Order, error) {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	var order order_model.Order
	idUser := ctx.Value(ctxkey.UserIDKey)
	err := db.QueryRowContext(ctx, "SELECT id_user, addres, messages, statuse, total_amount, price_all_products FROM orders WHERE id = $1 AND id_user = $2", idOrder, idUser).Scan(&order.User_id, &order.Address, &order.Messages, &order.Status, &order.Total_amount, &order.Price_all_products)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return order, order_model.ErrOrderNotFound
		}
		log.Error("error db",
			zap.Error(err),
		)
		return order, err
	}
	return order, nil
}
func (r *OrderUserRepo) GetPriceAllProduct(ctx context.Context, idOrder int) (float64, error) {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	var price float64
	err := db.QueryRowContext(ctx, "SELECT price_all_products FROM orders WHERE id = $1", idOrder).Scan(&price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return price, order_model.ErrOrderNotFound
		}
		log.Error("error db",
			zap.Error(err),
		)
		return price, err
	}
	return price, nil
}
func (r *OrderUserRepo) GetAllOrdersItems(ctx context.Context, idOrder int) ([]order_model.Order_items, error) {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	var orderItems []order_model.Order_items
	rows, err := db.QueryContext(ctx, "SELECT id_skus, quantity, price_at_purchase FROM orders_items WHERE id_order = $1", idOrder)
	if err != nil {
		log.Error("error db",
			zap.Error(err),
		)
		return orderItems, err
	}
	defer rows.Close()
	for rows.Next() {
		var Items order_model.Order_items
		if err := rows.Scan(&Items.Id_skus, &Items.Quantity, &Items.Price_at_purchase); err != nil {
			log.Error("error rows db",
				zap.Error(err),
			)
			return orderItems, err
		}
		orderItems = append(orderItems, Items)
	}
	if err := rows.Err(); err != nil {
		log.Error("error rows db",
			zap.Error(err),
		)
		return orderItems, err
	}
	return orderItems, nil
}
func (r *OrderUserRepo) FinalUpdate(ctx context.Context, total int, allPrice float64, status string, idOrder int) (int, error) {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	var newTotal int
	err := db.QueryRowContext(ctx, "UPDATE orders SET total_amount = $1,  price_all_products = $2, statuse = $3 WHERE id = $4 RETURNING total_amount", total, allPrice, status, idOrder).Scan(&newTotal)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return newTotal, order_model.ErrOrderNotFound
		}
		log.Error("error db",
			zap.Error(err),
		)
		return newTotal, err
	}
	return newTotal, nil
}
func (r *OrderUserRepo) AddOrderItems(ctx context.Context, idOrder, idSkus, quantity int, price_at_purchase float64) error {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	_, err := db.ExecContext(ctx, "INSERT INTO orders_items (id_order, id_skus, quantity, price_at_purchase) VALUES ($1, $2, $3, $4)", idOrder, idSkus, quantity, price_at_purchase)
	if err != nil {
		log.Error("error db",
			zap.Error(err),
		)
		return err
	}
	return nil
}
func (r *OrderUserRepo) GetUrl(ctx context.Context, idOrder int) (string, error) {
	log := logger.FromContext(ctx)
	var url string
	db := dto.Getter(ctx, r.db)
	idUser := ctx.Value(ctxkey.UserIDKey)
	err := db.QueryRowContext(ctx, "SELECT temporary_url_bank FROM orders WHERE id = $1 AND id_user = $2", idOrder, idUser).Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return url, order_model.ErrOrderNotFound
		}
		log.Error("error db",
			zap.Error(err),
		)
		return url, err
	}
	return url, nil
}
func (r *OrderUserRepo) AddOrUpdatedUrl(ctx context.Context, url string, idOrder int) error {
	log := logger.FromContext(ctx)
	var id string
	db := dto.Getter(ctx, r.db)
	idUser := ctx.Value(ctxkey.UserIDKey)
	err := db.QueryRowContext(ctx, "UPDATE orders SET temporary_url_bank = $1 WHERE id = $2 AND id_user = $3 RETURNING id", url, idOrder, idUser).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return order_model.ErrOrderNotFound
		}
		log.Error("error db",
			zap.Error(err),
		)
		return err
	}
	return nil
}
func (r *OrderUserRepo) AddUpdate_atOrder(ctx context.Context, idOrder int) (time.Time, error) {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	var timeUpdate time.Time
	err := db.QueryRowContext(ctx, "UPDATE orders SET updated_at = NOW() WHERE id = $1 RETURNING updated_at", idOrder).Scan(&timeUpdate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return timeUpdate, order_model.ErrOrderNotFound
		}
		log.Error("error db",
			zap.Error(err),
		)
		return timeUpdate, err
	}
	return timeUpdate, nil
}
func (r OrderUserRepo) GetBankPaymentId(ctx context.Context, idOrder int) (string, error) {
	log := logger.FromContext(ctx)
	var idPayment string
	db := dto.Getter(ctx, r.db)
	err := db.QueryRowContext(ctx, "SELECT bank_payment_id FROM orders WHERE id = $1", idOrder).Scan(&idPayment)
	if err != nil {
		log.Error("error db",
			zap.Error(err),
		)
		return idPayment, err
	}
	return idPayment, nil
}
func (r *OrderUserRepo) GetIdOrderFromByBankId(ctx context.Context, bankPaymentId string) (int, error) {
	log := logger.FromContext(ctx)
	var idOrder int
	db := dto.Getter(ctx, r.db)
	err := db.QueryRowContext(ctx, "SELECT id FROM orders WHERE bank_payment_id = $1", bankPaymentId).Scan(&idOrder)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, order_model.ErrOrderNotFound
		}
		log.Error("error db",
			zap.Error(err),
		)
		return 0, err
	}
	return idOrder, nil
}
func (r *OrderUserRepo) GetStatuseOrder(ctx context.Context, idOrder int) (string, error) {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	var statuse string
	err := db.QueryRowContext(ctx, "SELECT statuse FROM orders WHERE id = $1", idOrder).Scan(&statuse)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return statuse, order_model.ErrOrderNotFound
		}
		log.Error("error db",
			zap.Error(err),
		)
		return statuse, err
	}
	return statuse, nil
}
func (r *OrderUserRepo) UpdateStatusOrder(ctx context.Context, idOrder int, status string) (string, error) {
	log := logger.FromContext(ctx)
	db := dto.Getter(ctx, r.db)
	var newStatus string
	err := db.QueryRowContext(ctx, "UPDATE orders SET statuse = $1 WHERE id = $2 RETURNING statuse", status, idOrder).Scan(&newStatus)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", order_model.ErrOrderNotFound
		}
		log.Error("error db",
			zap.Error(err),
		)
		return "", err
	}
	return newStatus, nil
}

package cart_repo

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	cart_model "github.com/patato8984/Shop/internal/modules/cart/model"
)

type CartRepo struct {
	db *sql.DB
}

func NewCartRepo(db *sql.DB) *CartRepo {
	return &CartRepo{db: db}
}
func (r CartRepo) SearchCart(idUser int) (int, error) {
	var id int
	err := r.db.QueryRow("SELECT id FROM cart WHERE id_user = $1 AND status = 'active'", idUser).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return id, nil
		}
		return id, fmt.Errorf("error db: %s", err)
	}
	return id, nil
}
func (r CartRepo) CreatedCart(idUser int) (cart_model.Cart, error) {
	var cart cart_model.Cart
	err := r.db.QueryRow("INSERT INTO cart (id_user, created_at) VALUES ($1, NOW()) RETURNING created_at", idUser).Scan(&cart.Created_at)
	if err != nil {
		return cart, fmt.Errorf("error db: %s", err)
	}
	cart.Id_user = idUser
	return cart, nil
}
func (r CartRepo) GetIdCartItems(idCart int) (int, error) {
	var idCartItems int
	if err := r.db.QueryRow("SELECT id FROM cart_items WHERE id_cart = $1 LIMIT 1", idCart).Scan(&idCartItems); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("cart items not Faund")
		}
		return 0, fmt.Errorf("error db: %s", err)
	}
	return idCartItems, nil

}
func (r CartRepo) GetCart(idItems int) (cart_model.Cart, error) {
	rows, err := r.db.Query("SELECT c.id, c.id_user, c.created_at, c.update_at, c.status, i.id, i.id_cart, i.id_skus, i.quantity, i.price_snapshot FROM cart c JOIN cart_items i ON c.id = i.id_cart WHERE id = $1", idItems)
	var cart cart_model.Cart
	if err != nil {
		return cart, fmt.Errorf("error db: %s", err)
	}
	defer rows.Close()
	cart.Cart_items = []cart_model.Items{}
	for rows.Next() {
		items := cart_model.Items{}
		var idcart int
		var iduser int
		var created_at time.Time
		var update_at time.Time
		var status string
		if err := rows.Scan(&idcart, &iduser, &created_at, &update_at, &status, &items.Id, &items.Id_cart, &items.Id_skus, &items.Quantity, &items.Price_snapshot); err != nil {
			return cart, fmt.Errorf("error db: %s", err)
		}
		cart.Id_user = iduser
		cart.Created_at = created_at
		cart.Update_at = update_at
		cart.Status = status
		cart.Cart_items = append(cart.Cart_items, items)
	}
	return cart, nil
}
func (r CartRepo) AddUpdate_at(idUser int) (time.Time, error) {
	var time time.Time
	err := r.db.QueryRow("UPDATE cart SET update_at = NOW() WHERE id_user = $1 RETURNING update_at", idUser).Scan(&time)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return time, errors.New("cart not found")
		}
		return time, fmt.Errorf("error db: %s", err)
	}
	return time, nil
}
func (r CartRepo) AddProduct(idCart, idUser, idSkus, quantity int) (int, error) {
	var idCartItems int
	err := r.db.QueryRow("INSERT INTO cart_items (id_cart, id_skus, quantity) VALUES ($1, $2, $3) RETURNING id", idCart, idSkus, quantity).Scan(&idCartItems)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("cart_items not found")
		}
		return 0, fmt.Errorf("error db: %s", err)
	}
	return idCartItems, nil
}

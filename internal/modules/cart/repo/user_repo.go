package cart_repo

import (
	"database/sql"
	"fmt"
	"time"

	cart_model "github.com/patato8984/Shop/internal/modules/cart/model"
)

type CartRepo struct {
	db *sql.DB
}

func NewCatalogRepo(db *sql.DB) *CartRepo {
	return &CartRepo{db: db}
}
func (r CartRepo) SearchCart(idUser int) (bool, error) {
	var ok bool
	err := r.db.QueryRow("SELECT 1 FROM cart WHERE id_user = $1", idUser).Scan(&ok)
	if err != nil {
		return false, fmt.Errorf("error db: %s", err)
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return true, nil
}
func (r CartRepo) CreatedCart(idUser int) (cart_model.Cart, error) {
	var cart cart_model.Cart
	err := r.db.QueryRow("INSERT INTO cart (id_user) VALUES ($1) RETURNING created_at").Scan(&cart.Created_at)
	if err != nil {
		return cart, fmt.Errorf("error db: %s", err)
	}
	cart.Id_user = idUser
	return cart, nil
}

func (r CartRepo) GetCart(idUser int) (cart_model.Cart, error) {
	rows, err := r.db.Query("SELECT c.id_user, c.created_at, c.update_at, c.status, i.id, i.id_cart, i.id_skus, i.quantity, i.price_snapshot FROM cart c JOIN cart_items ON c.id_user = i.id_cart WHERE c.id_user = $1", idUser)
	var cart cart_model.Cart
	if err != nil {
		return cart, fmt.Errorf("error db: %s", err)
	}
	defer rows.Close()
	cart.Cart_items = []cart_model.Items{}
	for rows.Next() {
		items := cart_model.Items{}
		var iduser int
		var created_at time.Time
		var update_at time.Time
		var status string
		if err := rows.Scan(&iduser, &created_at, &update_at, &status, &items.Id, items.Id_cart, &items.Id_skus, &items.Quantity, &items.Price_snapshot); err != nil {
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
func (r CartRepo) SaveItems(items cart_model.Items) error {

}

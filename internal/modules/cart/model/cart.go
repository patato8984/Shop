package cart_model

import (
	"database/sql"
	"errors"
	"time"
)

type RequestDeltaAndIdSkus struct {
	IdSkus int `json:"id_skus"`
	Delta  int `json:"delta"`
}
type Cart struct {
	Id_cart             int          `json:"id_cart,omitempty"`
	Id_user             int          `json:"id_user,omitempty"`
	Created_at          time.Time    `json:"created_at,omitempty"`
	Update_at           sql.NullTime `json:"update_at,omitempty"`
	Status              string       `json:"status,omitempty"`
	PriceAllProductCart float64      `json:"priceAllProductCart,omitempty"`
	Cart_items          []Items      `json:"Cart_items,omitempty"`
}
type Items struct {
	Id             int     `jsonL:"id,omitempty"`
	Id_cart        int     `json:"id_cart,omitempty"`
	Id_skus        int     `json:"id_skus,omitempty"`
	Quantity       int     `json:"quantity,omitempty"`
	Price_snapshot float64 `json:"price_snapshot,omitempty"`
}
type CartEvent struct {
	EventType string `json:"type"`
	PayLoad   any    `json:"payload"`
}

type CartUpdate struct {
	IdCart int `json:"cart_id"`
}

var (
	ErrJson             = errors.New("error json")
	ErrUrl              = errors.New("incorrect URL")
	ErrCartNotFound     = errors.New("cart not found")
	ErrCartItemsNotFund = errors.New("cart items not found")
	ErrMaxQuantity      = errors.New("quantity max 50")
	ErrQuantityCatalog  = errors.New("there are not so many items in stock")
	ErrMinQuantity      = errors.New("quantity min 1")
)

func (c *Items) ChangeQuantityCatalogAndGetPrice(quantity int) (float64, error) {
	if c.Quantity < quantity {
		return 0, ErrQuantityCatalog
	}
	c.Price_snapshot *= float64(quantity)
	return c.Price_snapshot, nil
}
func (c *Items) ChangeQuantity(delta int) error {
	finalQuantity := c.Quantity + delta
	if finalQuantity > 50 {
		return ErrMaxQuantity
	} else if finalQuantity <= 0 {
		return ErrMinQuantity
	}
	c.Quantity = finalQuantity
	return nil
}

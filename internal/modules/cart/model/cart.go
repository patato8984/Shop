package cart_model

import "time"

type Cart = struct {
	Id_user    int       `json:"id_user,omitempty"`
	Created_at time.Time `json:"created_at,omitempty"`
	Update_at  time.Time `json:"update_at,omitempty"`
	Status     string    `json:"status,omitempty"`
	Cart_items []Items   `json:"Cart_items,omitempty"`
}
type Items = struct {
	Id             int     `jsonL:"id,omitempty"`
	Id_cart        int     `json:"id_cart,omitempty"`
	Id_skus        int     `json:"id_skus,omitempty"`
	Quantity       int     `json:"quantity,omitempty"`
	Price_snapshot float64 `json:"price_snapshot,omitempty"`
}

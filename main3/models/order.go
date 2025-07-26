package models

type Order = struct {
	Product Carts
	IdCart  int    `json:"idcart"`
	City    string `json:"city"`
	Address string `json:"address"`
}

package model

type Product = struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	SKUs []SKU  `json:"skus,omitempty"`
}
type SKU = struct {
	Id         int     `json:"id"`
	Product_id string  `json:"product_id"`
	Storage    int     `json:"storage,omitempty"`
	Colour     int     `json:"colour,omitempty"`
	Price      float32 `json:"price"`
	Stock      int     `json:"stock"`
}

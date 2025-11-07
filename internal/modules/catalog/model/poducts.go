package model

type Product = struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	SKUs []SKU  `json:"skus,omitempty"`
}
type SKU = struct {
	Id         int     `json:"id"`
	Product_id string  `json:"product_id"`
	Storage    int     `json:"storage"`
	Colour     int     `json:"colour"`
	Price      float32 `json:"price"`
	Stock      int     `json:"stock"`
}

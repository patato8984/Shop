package model

type Product = struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	SKUs []SKU
}
type SKU = struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	Storage int     `json:"storage"`
	Colour  int     `json:"colour"`
	Price   float32 `json:"price"`
	Stock   int     `json:"stock"`
}

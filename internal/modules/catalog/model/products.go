package catalog_model

import (
	"errors"
)

type Product = struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	SKUs []SKU  `json:"skus,omitempty"`
}
type SKU = struct {
	Id         int     `json:"id"`
	Product_id int     `json:"product_id"`
	Storage    int     `json:"storage,omitempty"`
	Colour     string  `json:"colour,omitempty"`
	Price      float64 `json:"price"`
	Stock      int     `json:"stock"`
}

type CatalogEvent struct {
	EventType string   `json:"type"`
	PayLoad   any      `json:"payload"`
	MetaDate  MetaDate `json:"meta_date"`
}
type MetaDate struct {
	IDUser   int    `json:"id_user"`
	RoleUser string `json:"role"`
}
type StockUpdatedLoad struct {
	ProductID int `json:"product_id"`
	SkusID    int `json:"skus_id"`
	NewStock  int `json:"new_stock"`
}
type PriceUpdatedLoad struct {
	ProductID int     `json:"product_id"`
	SkusID    int     `json:"skus_id"`
	NewPrice  float64 `json:"new_Price"`
}
type SkusDeleteLoad struct {
	ProductID int `json:"product_id"`
	SkusID    int `json:"skus_id"`
}

var (
	ErrJson            = errors.New("error json")
	ErrUrl             = errors.New("incorrect URL")
	ErrProductNotFound = errors.New("product not found")
	ErrShortID         = errors.New("short id")
	ErrSkuNotFound     = errors.New("sku not found")
	ErrShortName       = errors.New("short name")
)

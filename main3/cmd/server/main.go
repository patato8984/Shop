package main

import (
	"main3/internal/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("/products", handlers.Product)
	http.HandleFunc("/products/", handlers.Idproduct)
	http.HandleFunc("/cart", handlers.Carts)
	http.HandleFunc("/orders", handlers.OrderRequest)
	http.ListenAndServe(":8080", nil)
}

package handlers

import (
	"encoding/json"
	"main3/internal/storage"
	"main3/models"
	"net/http"
)

func Carts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(storage.BaseCarts)
	case http.MethodPost:
		var Reques struct {
			UserId  int    `json:"userid"`
			Product string `json:"Product"`
		}
		if err := json.NewDecoder(r.Body).Decode(&Reques); err != nil {
			http.Error(w, `{"json":"error"}`, http.StatusBadRequest)
		}
		if _, ok := storage.BaseProducts[Reques.Product]; !ok {
			http.Error(w, `{"json":"There is no product"}`, http.StatusNotFound)
			return
		}
		cart, result := storage.BaseCarts[Reques.UserId]
		if !result {
			cart = models.Carts{
				List: []string{},
				Sum:  0,
			}
		} else {
			cart = storage.BaseCarts[Reques.UserId]
		}
		cart.List = append(cart.List, Reques.Product)
		cart.Sum += storage.BaseProducts[Reques.Product].Price
		storage.BaseCarts[Reques.UserId] = cart
		w.Write([]byte(`{"status" : "Product added to cart"}`))
	}
}

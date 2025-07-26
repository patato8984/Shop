package handlers

import (
	"encoding/json"
	"main3/internal/storage"
	"main3/models"
	"net/http"
)

func OrderRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(storage.BaseOrder)
	case http.MethodPost:
		var BaseOrder struct {
			IdCart  int    `json:"idcart"`
			City    string `json:"city"`
			Address string `json:"address"`
		}
		if err := json.NewDecoder(r.Body).Decode(&BaseOrder); err != nil {
			http.Error(w, `{"status":"error"}`, http.StatusBadRequest)
		}
		if _, ok := storage.BaseCarts[BaseOrder.IdCart]; !ok {
			http.Error(w, `{"status":"not found carts"}`, http.StatusNotFound)
		}
		if len(storage.BaseCarts[BaseOrder.IdCart].List) == 0 {
			http.Error(w, `{"error": "Cart is empty"}`, http.StatusBadRequest)
			return
		}
		var order models.Order = models.Order{
			Product: storage.BaseCarts[BaseOrder.IdCart],
			IdCart:  BaseOrder.IdCart,
			City:    BaseOrder.City,
			Address: BaseOrder.Address,
		}
		storage.BaseOrder[order.IdCart] = order
		for _, base := range order.Product.List {
			for i, _ := range storage.BaseProducts {
				if base == i {
					delete(storage.BaseProducts, base)
				}
			}
		}
		delete(storage.BaseCarts, BaseOrder.IdCart)
		w.Write([]byte(`{"status":"purchase completed"}`))
	}
}

package handlers

import (
	"encoding/json"
	"main3/internal/storage"
	"main3/models"
	"net/http"
	"strconv"
	"strings"
)

func Idproduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/products/")
	ok2, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"status":"error"}`, http.StatusBadRequest)
		return
	}
	for kye, _ := range storage.BaseProducts {
		if storage.BaseProducts[kye].Id == ok2 {
			json.NewEncoder(w).Encode(storage.BaseProducts[kye])
		}
	}
}
func Product(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(storage.BaseProducts)
	case http.MethodPost:
		var Productes = models.Product{}

		if err := json.NewDecoder(r.Body).Decode(&Productes); err != nil {
			http.Error(w, `{"status" : "error"}`, http.StatusBadRequest)
		}
		storage.BaseProducts[Productes.Name] = Productes
		w.Write([]byte(`{"status" : "complete"}`))
	}
}

package user

import (
	"encoding/json"
	"net/http"

	"github.com/patato8984/Shop/internal/modules/user/model"
	"github.com/patato8984/Shop/internal/shared/dto"
)

type UserHandler struct {
	service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "invalid json"})
		return
	}
	err := h.service.RegisterCase(user)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		switch err.Error() {
		case "short password or Nickname":
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.Response(err.Error(), http.StatusBadRequest, nil))
		case "nickname busy":
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(dto.Response(err.Error(), http.StatusConflict, nil))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.Response("error server", http.StatusInternalServerError, nil))
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
func (h *UserHandler) Authentication(w http.ResponseWriter, r http.Request) {
	var passwordAndName model.User
	if err := json.NewDecoder(r.Body).Decode(&passwordAndName); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Response("error json", http.StatusBadRequest, nil))
		return
	}
	user, err := h.service.GetToken(passwordAndName.Nickname, passwordAndName.Password)
	if err != nil {
		switch err.Error() {
		case "user not found":
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.Response(err.Error(), http.StatusNotFound, nil))
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.Response("error server", http.StatusInternalServerError, nil))
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.Response("seccus", http.StatusOK, user))
}

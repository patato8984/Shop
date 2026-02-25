package dto

import (
	"encoding/json"
	"log"
	"net/http"
)

type MessageResponse = struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    any    `json:"data,omitempty"`
}

func Response(w http.ResponseWriter, message string, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(MessageResponse{Message: message, Status: status, Data: data})
	if err != nil {
		log.Printf("critical: could not send error response: %v", err)
	}
}

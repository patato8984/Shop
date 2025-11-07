package dto

type MessageResponse = struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    any    `json:"data,omitempty"`
}

func Response(message string, status int, data any) MessageResponse {
	return MessageResponse{Message: message, Status: status, Data: data}
}

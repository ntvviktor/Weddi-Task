package utils

import (
	"net/http"
	"encoding/json"
)

type ErrorResponse struct {
	Message string `json:"message"`
}


func NewErrorResponse(w http.ResponseWriter, statusCode int, response string) {
	error := ErrorResponse{
		response,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(&error)
}
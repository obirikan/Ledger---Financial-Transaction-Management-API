// internal/response/json_responder.go

package response

import (
	"encoding/json"
	"log"
	"net/http"
)

type Responder interface {
	Error(w http.ResponseWriter, statusCode int, message string)
	Success(w http.ResponseWriter, data interface{}, message string)
}

type JSONResponder struct{}

func NewJSONResponder() Responder {
	return &JSONResponder{}
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func (r *JSONResponder) Error(w http.ResponseWriter, statusCode int, message string) {
	resp := ErrorResponse{
		Error:   "API Error",
		Code:    statusCode,
		Message: message,
	}
	r.sendJSON(w, statusCode, resp)
}

func (r *JSONResponder) Success(w http.ResponseWriter, data interface{}, message string) {
	resp := SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	}
	r.sendJSON(w, http.StatusOK, resp)
}

func (r *JSONResponder) sendJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Error encoding JSON: %v", err)
	}
}



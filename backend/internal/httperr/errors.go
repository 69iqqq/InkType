package httperr

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func JSON(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error: msg,
		Code:  code,
	})
}

func InternalError(w http.ResponseWriter) {
	JSON(w, "internal server error", http.StatusInternalServerError)
}

func BadRequest(w http.ResponseWriter, msg string) {
	JSON(w, msg, http.StatusBadRequest)
}

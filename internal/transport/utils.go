package transport

import (
	"encoding/json"
	"log"
	"net/http"
)

type ApiError struct {
	Err            error `json:"err,omitempty"`
	HttpStatusCode int   `json:"http_status_code,omitempty"`
}

func WriteError(w http.ResponseWriter, apiError ApiError) {
	log.Printf("Error occurred: %v", apiError.Err)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(apiError.HttpStatusCode)

	json.NewEncoder(w).Encode(apiError)
}

func WriteJson(w http.ResponseWriter, statusCode int, val any) {
	log.Printf("Writing JSON response with status code: %d", statusCode)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if val == nil {
		w.Write([]byte{})
	}

	json.NewEncoder(w).Encode(val)
}

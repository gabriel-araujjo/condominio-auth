package errors

import (
	"encoding/json"
	"net/http"
)

type errorMessage struct {
	Message string `json:"error"`
}

// WriteErrorWithCode sends an error back with the specified status code
func WriteErrorWithCode(w http.ResponseWriter, status int, errToSend interface{}) error {
	w.WriteHeader(status)
	if message, ok := errToSend.(string); ok {
		return json.NewEncoder(w).Encode(&errorMessage{message})
	}
	return json.NewEncoder(w).Encode(errToSend)
}

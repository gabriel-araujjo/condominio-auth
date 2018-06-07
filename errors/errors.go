package errors

import (
	"encoding/json"
	"net/http"
)

type errorMessage struct {
	Message string `json:"message"`
}

// WriteErrorWithCode sends an error back with the specified status code
func WriteErrorWithCode(w http.ResponseWriter, status int, errToSend interface{}) (int, error) {
	w.WriteHeader(status)
	if message, ok := errToSend.(string); ok {
		return writeError(w, &errorMessage{message})
	}
	return writeError(w, errToSend)
}

func writeError(w http.ResponseWriter, errToSend interface{}) (int, error) {
	data, err := json.Marshal(errToSend)
	if err != nil {
		return 0, err
	}
	w.Header().Set("Content-Type", "application/json")
	return w.Write(data)
}

package mock

import (
	"net/http"
)

type MockResponseWriter struct {
	Head   http.Header
	Body   []byte
	Status int
}

func (w *MockResponseWriter) Header() http.Header {
	return w.Head
}

func (w *MockResponseWriter) Write(body []byte) (int, error) {
	w.Body = make([]byte, len(body))
	copy(w.Body, body)
	return len(w.Body), nil
}

func (w *MockResponseWriter) WriteHeader(status int) {
	w.Status = status
}

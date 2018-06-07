package routes

import (
	"fmt"
	"net/http"

	"github.com/gabriel-araujjo/condominio-auth/errors"
)

func checkContentType(desiredType string) *Middleware {
	return newMiddleware(func(resp http.ResponseWriter, req *http.Request) bool {
		contentType := req.Header.Get("Content-Type")
		if contentType != desiredType {
			err := struct {
				Message  string
				Expected string
			}{
				Message:  fmt.Sprintf("unsupported \"Content-Type: %s\"", contentType),
				Expected: desiredType,
			}
			errors.WriteErrorWithCode(resp, http.StatusUnsupportedMediaType, err)
			return true
		}
		return false
	})
}

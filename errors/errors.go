package httperrors

import (
	"github.com/gin-gonic/gin"
	"errors"
	"net/http"
)

const (
	ErrNotFound gin.ErrorType = 1 << 50 + http.StatusNotFound
	ErrForbidden gin.ErrorType = 1 << 50 + http.StatusForbidden
	ErrPreconditionFailed gin.ErrorType = 1 << 50 + http.StatusPreconditionFailed
)

func NotFound(message string) *gin.Error {
	return err(ErrNotFound, message)
}

func Forbidden(message string) *gin.Error {
	return err(ErrForbidden, message)
}

func PreconditionFailed(message string) *gin.Error {
	return err(ErrPreconditionFailed, message)
}

func err(code gin.ErrorType, msg string) *gin.Error {
	return &gin.Error{
		Type: code,
		Err: errors.New(msg),
	}
}

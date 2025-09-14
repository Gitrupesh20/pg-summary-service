package domain

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"net/http"
	"pg-summary-service/internal/logger"
)

var (
	ErrInvalidInput               = errors.New("invalid input")
	ErrExternalServiceUnreachable = errors.New("external service unreachable")
)

// AppError is a custom error with a message and status code
type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewNotFoundError(msg string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: msg}
}

func NewNoContentError(msg string) *AppError {
	return &AppError{Code: http.StatusNoContent, Message: msg}
}

func NewBadRequestError(msg string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: msg}
}

func NewInternalError(msg string) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: msg}
}

// HandlePGError converts DB errors into AppError
func HandlePGError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return NewNotFoundError("no data found")
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "42P01":
			return NewNotFoundError("requested table or list does not exist")
		case "23505":
			return NewBadRequestError("duplicate entry, already exists")
		case "23503":
			logger.Log.Error("foreign key violation", zap.Error(err))
			return NewInternalError("internal server error")
		default:
			logger.Log.Error("postgres error", zap.String("code", string(pqErr.Code)), zap.Error(err))
			return NewInternalError("internal server error")
		}
	}

	logger.Log.Error("unexpected DB error", zap.Error(err))
	return NewInternalError("internal server error")
}

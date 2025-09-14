package utils

import (
	"pg-summary-service/internal/domain"
)

func ValidateDBDetails(data domain.RemoteDBDetails) error {

	if data.Host == "" {
		return domain.NewBadRequestError("host cannot be empty")
	}
	if data.DBName == "" {
		return domain.NewBadRequestError("database name cannot be empty")
	}
	if data.Password == "" {
		return domain.NewBadRequestError("password cannot be empty")
	}
	if data.Port == 0 {
		return domain.NewBadRequestError("port cannot be empty")
	}
	if data.User == "" {
		return domain.NewBadRequestError("user cannot be empty")
	}
	return nil
}

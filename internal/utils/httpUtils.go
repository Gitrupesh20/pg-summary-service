package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"pg-summary-service/internal/domain"
	logger1 "pg-summary-service/internal/logger"
	"strconv"
	"strings"
	"time"
)

const (
	MaxRetry     = 30
	MaxSleepTime = 30 * time.Second
)

func PostWithRetry(url string, noOfRetry int, payload any) (*http.Response, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	if noOfRetry < 1 {
		noOfRetry = 1
	} else if noOfRetry >= MaxRetry {
		noOfRetry = MaxRetry // prevent too many retries
	}

	for try := 0; try < noOfRetry; try++ {
		resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
		if err != nil {
			logger1.Log.Warn(fmt.Sprintf("request failed (try %d), retrying...", try+1), zap.Error(err))
		} else {
			// check status code
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return resp, nil
			}

			// client error, don’t retry
			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				resp.Body.Close()
				return nil, domain.NewBadRequestError(fmt.Sprintf("external API returned %d", resp.StatusCode))
			}
			// server error retry
			resp.Body.Close()
			logger1.Log.Warn(fmt.Sprintf("server error (status %d), retrying (try %d)...", resp.StatusCode, try+1))
		}

		// exponential backoff with max cap
		sleepTime := time.Second * time.Duration(1<<try)
		if sleepTime > MaxSleepTime {
			sleepTime = MaxSleepTime
		}
		time.Sleep(sleepTime)
	}

	return nil, domain.ErrExternalServiceUnreachable
}

func ExtractIDFromPath(r *http.Request, prefix string) (string, error) {
	path := r.URL.Path
	if !strings.HasPrefix(path, prefix) {
		return "", errors.New("invalid path")
	}
	id := strings.TrimPrefix(path, prefix)
	id = strings.Trim(id, "/")
	if id == "" {
		return "", errors.New("id not provided in path")
	}
	return id, nil
}

func ParseQueryInt(r *http.Request, key string, defaultValue int) int {
	values := r.URL.Query()
	if raw := values.Get(key); raw != "" {
		if val, err := strconv.Atoi(raw); err == nil {
			return val
		}
	}
	return defaultValue
}
func SendError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	// If it's an AppError
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		http.Error(w, appErr.Message, appErr.Code)
		return
	}

	http.Error(w, "internal server error", http.StatusInternalServerError)
}

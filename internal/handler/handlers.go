package handler

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"pg-summary-service/internal/domain"
	logger1 "pg-summary-service/internal/logger"
	"pg-summary-service/internal/utils"
)

const (
	defaultLimit  = 0
	defaultOffset = 20
)

func SyncSummaryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse request body
	var req domain.RemoteDBDetails
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger1.Log.Error("Error decoding request", zap.Error(err))
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	} else if err = utils.ValidateDBDetails(req); err != nil {
		utils.SendError(w, err)
	}

	_, err := service.SyncSummary(r.Context(), req)
	if err != nil {
		logger1.Log.Error("error while fetching and saving data from external api", zap.Error(err))
		utils.SendError(w, err)
		return
	}

	// Placeholder response
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "sync triggered",
	})
}

func GetSummariesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("geckjnn")

	// Parse optional query params for pagination
	//note: limit and offset to control fetching to avoid over-fetching too many rows at once
	offset := utils.ParseQueryInt(r, "limit", defaultLimit)
	limit := utils.ParseQueryInt(r, "offset", defaultOffset) //default limit 20

	if resp, err := service.GetSummaries(r.Context(), offset, limit); err != nil {
		logger1.Log.Error("error at GetSummaries handler", zap.Error(err))
		utils.SendError(w, err)
		return
	} else {
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func GetSummaryByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract ID from URL using utils
	id, err := utils.ExtractIDFromPath(r, "/summaries/")
	if err != nil {
		http.Error(w, "missing or invalid summary ID", http.StatusBadRequest)
		return
	}

	if resp, err := service.GetSummaryByID(r.Context(), id); err != nil {
		logger1.Log.Error("error at GetSummaryByIDHandler handler", zap.Error(err))
		utils.SendError(w, err)
		return
	} else {
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>404 Not Found</title>
		<style>
			body { font-family: Arial, sans-serif; text-align: center; padding: 50px; }
			h1 { font-size: 50px; }
			p { font-size: 20px; }
		</style>
	</head>
	<body>
		<h1>404</h1>
		<p>Oops! The page you requested does not exist.</p>
	</body>
	</html>
	`
	w.Write([]byte(html))
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"pg-summary-service/internal/logger"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"pg-summary-service/internal/config"
	"pg-summary-service/internal/handler"
	"pg-summary-service/internal/repository/external"
	"pg-summary-service/internal/repository/local"
	"pg-summary-service/internal/service"
)

func main() {
	// Load config
	if err := config.LoadConfig(); err != nil {
		log.Fatal("failed to load config: %w", err)
	}

	logger.Init(config.Debug(), config.GetLogDir(), config.GetLogFile()) // false = prod (JSON logs), true = dev (console logs)
	defer logger.Log.Sync()

	logger.Log.Info("Starting server...")

	if err := startServer(); err != nil {
		logger.Log.Fatal("Server failed", zap.Error(err))
	}
}

func startServer() error {

	// created db pool and Initialize tables
	pool, err := initDB()
	if err != nil {
		return err
	}
	defer pool.Close()

	// Instantiate Repositories
	extRepo := external.NewExternalRepository(config.GetExternalDbUrl(), config.GetRetries())
	localRepo := local.NewLocalRepository(pool)

	// Service
	svc := service.NewSummaryService(extRepo, localRepo)

	// Register routes
	handler.RegisterRoutes(http.HandleFunc, *svc)

	// Start server
	port := fmt.Sprintf(":%s", config.GetPort())
	logger.Log.Info("Server starting", zap.String("port", port))
	fmt.Println("*************************************************| Starting server |*************************************************")

	if err := http.ListenAndServe(port, nil); err != nil {
		return fmt.Errorf("error starting server: %w", err)
	}

	return nil
}

func initDB() (*pgxpool.Pool, error) {
	// PostgreSQL connection pool configuration
	dbConfig, err := pgxpool.ParseConfig(config.GetLocalDbUrl())
	if err != nil {
		return nil, fmt.Errorf("failed to parse DB config: %w", err)
	}

	dbStats := config.GetDBStats()
	dbConfig.MaxConns = int32(dbStats.MaxConnections)
	dbConfig.MinConns = int32(dbStats.MaxIdleConnections)
	dbConfig.MaxConnLifetime = dbStats.MaxConnectionLifeTime

	// Connect to PostgreSQL
	pool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to local DB: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	queries := []string{
		`CREATE TABLE IF NOT EXISTS summaries (id VARCHAR PRIMARY KEY, source_info VARCHAR, synced_at TIMESTAMP)`,
		`CREATE TABLE IF NOT EXISTS schemas (id VARCHAR PRIMARY KEY, summary_id VARCHAR REFERENCES summaries(id), name VARCHAR)`,
		`CREATE TABLE IF NOT EXISTS tables (id VARCHAR PRIMARY KEY, schema_id VARCHAR REFERENCES schemas(id), name VARCHAR, row_count BIGINT, size_mb FLOAT)`,
	}

	for _, q := range queries {
		if _, err := pool.Exec(ctx, q); err != nil {
			return nil, fmt.Errorf("failed to init DB: %w", err)
		}
	}
	return pool, nil
}

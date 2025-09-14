package local

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"pg-summary-service/internal/domain"
	"pg-summary-service/internal/logger"
	"time"
)

type LocalRepository struct {
	db *pgxpool.Pool
}

func NewLocalRepository(db *pgxpool.Pool) *LocalRepository {
	return &LocalRepository{db: db}
}

func (lRepo *LocalRepository) AddSummary(ctx context.Context, src string, data *domain.ExternalSummaryResp) (any, error) {

	if src == "" {
		return nil, domain.NewBadRequestError("src cannot be an empty string")
	} else if data == nil {
		return nil, domain.NewBadRequestError("data cannot be an empty")
	}

	// Insert summary
	id := data.Id //uuid.New().String()
	syncedAt := time.Now()
	query := `INSERT INTO summaries (id, source_info, synced_at) VALUES ($1, $2, $3)`
	if _, err := lRepo.db.Exec(ctx, query, id, src, syncedAt); err != nil {
		logger.Log.Error("error while saving data to local db", zap.Error(err), zap.Any("data", data))
		return nil, domain.HandlePGError(err)
	}

	// Insert schemas and tables
	for _, schema := range data.Schemas {
		schemaID := uuid.New().String()
		query = `INSERT INTO schemas (id, summary_id, name) VALUES ($1, $2, $3)`
		_, err := lRepo.db.Exec(ctx, query, schemaID, id, schema.Name)
		if err != nil {
			logger.Log.Error("error while saving schema data to local db", zap.Error(err), zap.Any("schema", schema))
			return nil, domain.HandlePGError(err)
		}
		tableQuery := `INSERT INTO tables (id, schema_id, name, row_count, size_mb) VALUES ($1, $2, $3, $4, $5)`
		for _, table := range schema.Tables {
			_, err = lRepo.db.Exec(ctx, tableQuery,
				uuid.New().String(), schemaID, table.Name, table.TotalRows, table.Size)
			if err != nil {
				logger.Log.Error("error while saving table data to local db", zap.Error(err), zap.Any("table", table))
				return nil, domain.HandlePGError(err)
			}
		}
	}
	return nil, nil
	//return &domain.Summary{
	//	ID:        id,
	//	SourceInfo: sourceInfo,
	//	SyncedAt:  syncedAt,
	//	Schemas:   externalResp.Schemas,
	//}, nil
}

func (lRepo *LocalRepository) GetSummaryById(ctx context.Context, id string) (*domain.LocalSummaryByIdResp, error) {
	if id == "" {
		return nil, domain.NewBadRequestError("id cannot be an empty string")
	}

	query := `
	SELECT 
		s.id, s.source_info, s.synced_at,
		sc.id, sc.name,
		COUNT(t.id), 
		COALESCE(SUM(t.row_count), 0),
		COALESCE(SUM(t.size_mb), 0)
	FROM summaries s
	LEFT JOIN schemas sc ON sc.summary_id = s.id
	LEFT JOIN tables t ON t.schema_id = sc.id
	WHERE s.id = $1
	GROUP BY s.id, s.source_info, s.synced_at, sc.id, sc.name;
	`

	rows, err := lRepo.db.Query(ctx, query, id)
	if err != nil {
		logger.Log.Error("error while fetching summary", zap.Error(err), zap.Any("summary id", id))
		return nil, domain.HandlePGError(err)
	}
	defer rows.Close()

	var summary domain.LocalSummaryByIdResp
	firstRow := true

	for rows.Next() {
		var (
			summaryID   string
			source      string
			syncedAt    time.Time
			schemaID    *string
			schemaName  *string
			tableCount  int
			totalRows   int64
			totalSizeMb float64
		)

		if err = rows.Scan(&summaryID, &source, &syncedAt, &schemaID, &schemaName, &tableCount, &totalRows, &totalSizeMb); err != nil {
			logger.Log.Error("error while s-caning summary data", zap.Error(err))
			return nil, err
		}

		// fill top-level summary once
		if firstRow {
			summary.ID = summaryID
			summary.Source = source
			summary.SyncedAt = syncedAt
			firstRow = false
		}

		// add schema if exists
		if schemaID != nil && *schemaID != "" {
			schema := domain.SchemaSummary{
				Id:          *schemaID,
				Name:        *schemaName,
				TableCount:  tableCount,
				TotalRows:   totalRows,
				TotalSizeMB: totalSizeMb,
			}
			summary.Schemas = append(summary.Schemas, schema)
		}
	}

	if firstRow {
		return nil, domain.NewNotFoundError(fmt.Sprintf("summary with id %s not found", id))
	}

	return &summary, nil
}

func (lRepo *LocalRepository) GetSummary(ctx context.Context, offset int, limit int) ([]domain.LocalSummaryListItem, error) {

	query := `SELECT id, source_info, synced_at 
	          FROM summaries 
	          ORDER BY synced_at DESC 
	          LIMIT $1 OFFSET $2`
	rows, err := lRepo.db.Query(ctx, query, limit, offset)
	if err != nil {
		logger.Log.Error("error while fetching summaries", zap.Error(err))
		return nil, domain.HandlePGError(err)
	}
	defer rows.Close()

	var items []domain.LocalSummaryListItem
	for rows.Next() {
		var item domain.LocalSummaryListItem
		if err = rows.Scan(&item.ID, &item.DBName, &item.SyncedAt); err != nil { // Note: source_info as DBName for list
			logger.Log.Error("error while s-caning summaries", zap.Error(err))
			return nil, domain.HandlePGError(err)
		}
		items = append(items, item)
	}
	return items, nil
}

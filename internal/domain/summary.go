package domain

import "time"

/* ########################################## Local Struct ########################################## */

type LocalSummary struct {
	Id     int       `json:"summary_id"`
	Source string    `json:"source"`
	SyncAt time.Time `json:"sync_at"`
}
type LocalSchema struct {
	Id        string `json:"id"`
	SummaryId string `json:"summary_id"`
	Name      string `json:"name"`
}

type LocalTable struct {
	Id        string `json:"id"`
	SchemaId  string `json:"schema_id"`
	Name      string `json:"name"`
	TotalRows int    `json:"row_count"`
	Size      int    `json:"size_mb"`
}
type LocalSummaryListItem struct {
	ID       string    `json:"id"`
	DBName   string    `json:"db_name"`
	SyncedAt time.Time `json:"synced_at"`
}
type LocalSummaryByIdResp struct {
	ID       string          `json:"summary_id"`
	Source   string          `json:"source"`
	SyncedAt time.Time       `json:"synced_at"`
	Schemas  []SchemaSummary `json:"schemas"`
}

/* ########################################## External Struct ########################################## */

type ExternalSummaryResp struct {
	Id      string   `json:"summary_id"`
	Schemas []Schema `json:"schemas"`
}

/* ########################################## Common Struct ########################################## */

type Table struct {
	Name      string  `json:"name"`
	TotalRows int     `json:"row_count"`
	Size      float64 `json:"size_mb"`
}
type Schema struct {
	Name   string  `json:"name"`
	Tables []Table `json:"tables"`
}

type SchemaSummary struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	TableCount  int     `json:"table_count"`
	TotalRows   int64   `json:"total_rows"`
	TotalSizeMB float64 `json:"total_size_mb"`
}

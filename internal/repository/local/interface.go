package local

import (
	"context"
	"pg-summary-service/internal/domain"
)

type Local interface {
	GetSummaryById(ctx context.Context, id string) (*domain.LocalSummaryByIdResp, error)
	AddSummary(ctx context.Context, src string, data *domain.ExternalSummaryResp) (any, error)
	GetSummary(ctx context.Context, offset int, limit int) ([]domain.LocalSummaryListItem, error)
}

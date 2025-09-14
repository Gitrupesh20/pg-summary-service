package service

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"pg-summary-service/internal/domain"
	logger2 "pg-summary-service/internal/logger"
	"pg-summary-service/internal/repository/external"
	"pg-summary-service/internal/repository/local"
)

type SummaryService struct {
	externalRepo external.External
	localRepo    local.Local
}

func NewSummaryService(extRepo external.External, localRepo local.Local) *SummaryService {
	return &SummaryService{externalRepo: extRepo, localRepo: localRepo}
}

func (s *SummaryService) SyncSummary(ctx context.Context, details domain.RemoteDBDetails) (any, error) {
	if details.Host == "" || details.DBName == "" || details.User == "" || details.Password == "" || details.Port == 0 {
		return nil, domain.NewBadRequestError("invalid input")
	}

	externalResp, err := s.externalRepo.FetchSummaries(details)
	if err != nil {
		logger2.Log.Error("src :SyncSummary error while fetching from external repo: ", zap.Error(err))
		return nil, err
	}

	sourceInfo := fmt.Sprintf("%s:%s", details.Host, details.DBName) // Don't store pass
	return s.localRepo.AddSummary(ctx, sourceInfo, externalResp)
}

func (s *SummaryService) GetSummaries(ctx context.Context, offset, limit int) ([]domain.LocalSummaryListItem, error) {

	return s.localRepo.GetSummary(ctx, offset, limit)
}

func (s *SummaryService) GetSummaryByID(ctx context.Context, id string) (*domain.LocalSummaryByIdResp, error) {
	return s.localRepo.GetSummaryById(ctx, id)
}

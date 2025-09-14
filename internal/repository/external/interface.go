package external

import "pg-summary-service/internal/domain"

type External interface {
	FetchSummaries(details domain.RemoteDBDetails) (*domain.ExternalSummaryResp, error)
}

package external

import (
	"encoding/json"
	"fmt"
	"pg-summary-service/internal/domain"
	"pg-summary-service/internal/utils"
)

type ExternalRepository struct {
	URL     string
	Retries int
}

func NewExternalRepository(url string, retries int) *ExternalRepository {
	// constructor
	return &ExternalRepository{
		URL:     url,
		Retries: retries,
	}
}

func (eRepo *ExternalRepository) FetchSummaries(data domain.RemoteDBDetails) (*domain.ExternalSummaryResp, error) {
	resp, err := utils.PostWithRetry(eRepo.URL, eRepo.Retries, data)
	if err != nil {
		return nil, fmt.Errorf("error while fetching external summary list: %w", err)
	}
	defer resp.Body.Close()

	var externalSummaryResp domain.ExternalSummaryResp
	if err := json.NewDecoder(resp.Body).Decode(&externalSummaryResp); err != nil {
		return nil, domain.NewInternalError(
			fmt.Sprintf("failed to parse external summary list (status %d): %v", resp.StatusCode, err),
		)
	}
	// check for missing or empty schemas
	if externalSummaryResp.Id == "" || len(externalSummaryResp.Schemas) == 0 {
		return nil, domain.NewNotFoundError("external summary list not found or empty")
	}

	return &externalSummaryResp, nil
}

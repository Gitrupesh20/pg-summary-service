package test

import (
	"context"
	"errors"
	"os"
	"pg-summary-service/internal/domain"
	"pg-summary-service/internal/logger"
	"pg-summary-service/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Initialize logger for tests
func TestMain(m *testing.M) {
	//"./logs",
	//	logFile:               "server.log",
	logger.Init(true, "./logs", "server.log") // debug mode
	os.Exit(m.Run())
}

// Mock External Repository
type MockExtRepo struct {
	mock.Mock
}

func (m *MockExtRepo) FetchSummaries(details domain.RemoteDBDetails) (*domain.ExternalSummaryResp, error) {
	args := m.Called(details)
	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}
	return resp.(*domain.ExternalSummaryResp), args.Error(1)
}

// Mock Local Repository
type MockLocalRepo struct {
	mock.Mock
}

func (m *MockLocalRepo) AddSummary(ctx context.Context, sourceInfo string, resp *domain.ExternalSummaryResp) (any, error) {
	args := m.Called(ctx, sourceInfo, resp)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*domain.LocalSummaryByIdResp), args.Error(1)
}

func (m *MockLocalRepo) GetSummary(ctx context.Context, offset, limit int) ([]domain.LocalSummaryListItem, error) {
	args := m.Called(ctx, offset, limit)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.([]domain.LocalSummaryListItem), args.Error(1)
}

func (m *MockLocalRepo) GetSummaryById(ctx context.Context, id string) (*domain.LocalSummaryByIdResp, error) {
	args := m.Called(ctx, id)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*domain.LocalSummaryByIdResp), args.Error(1)
}

// Test SyncSummary Success
func TestSyncSummary(t *testing.T) {
	mockExt := new(MockExtRepo)
	mockLocal := new(MockLocalRepo)
	svc := service.NewSummaryService(mockExt, mockLocal)

	details := domain.RemoteDBDetails{
		Host:     "test",
		Port:     5432,
		User:     "user",
		Password: "pass",
		DBName:   "db",
	}

	extResp := &domain.ExternalSummaryResp{
		Id:      "summary1",
		Schemas: []domain.Schema{},
	}

	localResp := &domain.LocalSummaryByIdResp{
		ID: "local1",
	}

	mockExt.On("FetchSummaries", details).Return(extResp, nil)
	mockLocal.On("AddSummary", mock.Anything, "test:db", extResp).Return(localResp, nil)

	resAny, err := svc.SyncSummary(context.Background(), details)
	assert.NoError(t, err)

	// Type assert 'any' to expected type
	res := resAny.(*domain.LocalSummaryByIdResp)
	assert.Equal(t, "local1", res.ID)

	mockExt.AssertExpectations(t)
	mockLocal.AssertExpectations(t)
}

// Test SyncSummary Invalid Input
func TestSyncSummaryInvalidInput(t *testing.T) {
	mockExt := new(MockExtRepo)
	mockLocal := new(MockLocalRepo)
	svc := service.NewSummaryService(mockExt, mockLocal)

	details := domain.RemoteDBDetails{} // missing all fields

	res, err := svc.SyncSummary(context.Background(), details)
	assert.Nil(t, res)
	assert.EqualError(t, err, "invalid input") // compare error string
}

// Test External Repo Unreachable
func TestSyncSummaryExternalError(t *testing.T) {
	mockExt := new(MockExtRepo)
	mockLocal := new(MockLocalRepo)
	svc := service.NewSummaryService(mockExt, mockLocal)

	details := domain.RemoteDBDetails{
		Host:     "test",
		Port:     5432,
		User:     "user",
		Password: "pass",
		DBName:   "db",
	}

	mockExt.On("FetchSummaries", details).Return(nil, errors.New("external service down"))

	res, err := svc.SyncSummary(context.Background(), details)
	assert.Nil(t, res)
	assert.EqualError(t, err, "external service down")

	mockExt.AssertExpectations(t)
}

// Test GetSummaries
func TestGetSummaries(t *testing.T) {
	mockExt := new(MockExtRepo)
	mockLocal := new(MockLocalRepo)
	svc := service.NewSummaryService(mockExt, mockLocal)

	expected := []domain.LocalSummaryListItem{
		{ID: "1", DBName: "test:db"},
		{ID: "2", DBName: "test:db"},
	}

	mockLocal.On("GetSummary", mock.Anything, 0, 10).Return(expected, nil)

	res, err := svc.GetSummaries(context.Background(), 0, 10)

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
	mockLocal.AssertExpectations(t)
}

// Test GetSummaryByID
func TestGetSummaryByID(t *testing.T) {
	mockExt := new(MockExtRepo)
	mockLocal := new(MockLocalRepo)
	svc := service.NewSummaryService(mockExt, mockLocal)

	id := "local1"
	expected := &domain.LocalSummaryByIdResp{ID: id}

	mockLocal.On("GetSummaryById", mock.Anything, id).Return(expected, nil)

	res, err := svc.GetSummaryByID(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, id, res.ID)

	mockLocal.AssertExpectations(t)
}

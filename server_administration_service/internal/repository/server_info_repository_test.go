package repository_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"server_administration_service/internal/repository"
	"testing"
	"time"

	"github.com/elastic/go-elasticsearch/v9/esapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockESClient struct {
	mock.Mock
}

// Add Index method to satisfy ElasticsearchClient interface
func (m *MockESClient) Index(ctx context.Context, index string, body []byte) error {
	args := m.Called(ctx, index, body)
	return args.Error(0)
}

func (m *MockESClient) Search(ctx context.Context, index string, buf bytes.Buffer) (*esapi.Response, error) {
	args := m.Called(ctx, index, buf)
	if (args.Get(0) == nil) {
		return nil, args.Error(1)
	}
	return args.Get(0).(*esapi.Response), args.Error(1)
}

func TestGetNumServers_Success(t *testing.T) {
	gdb, mock, cleanup := repository.SetupMockDB(t)
	defer cleanup()

	repo := repository.NewServerInfoRepository(gdb, nil)

	mock.ExpectQuery("SELECT count\\(\\*\\) FROM \"servers\"").
		WillReturnRows(mock.NewRows([]string{"count"}).AddRow(5))

	num, err := repo.GetNumServers()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if num != 5 {
		t.Errorf("expected 5 servers, got %d", num)
	}
}

func TestGetNumServers_Error(t *testing.T) {
	gdb, mock, cleanup := repository.SetupMockDB(t)
	defer cleanup()

	repo := repository.NewServerInfoRepository(gdb, nil)

	mock.ExpectQuery("SELECT count\\(\\*\\) FROM \"servers\"").
		WillReturnError(fmt.Errorf("db error"))

	num, err := repo.GetNumServers()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if num != 0 {
		t.Errorf("expected 0 servers on error, got %d", num)
	}
}

func TestGetNumOnServers_Success(t *testing.T) {
	gdb, mock, cleanup := repository.SetupMockDB(t)
	defer cleanup()

	repo := repository.NewServerInfoRepository(gdb, nil)

	mock.ExpectQuery("SELECT count\\(\\*\\) FROM \"servers\" WHERE status = ?").
		WithArgs("On").
		WillReturnRows(mock.NewRows([]string{"count"}).AddRow(3))

	num, err := repo.GetNumOnServers()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if num != 3 {
		t.Errorf("expected 3 ON servers, got %d", num)
	}
}

func TestGetNumOnServers_Error(t *testing.T) {
	gdb, mock, cleanup := repository.SetupMockDB(t)
	defer cleanup()

	repo := repository.NewServerInfoRepository(gdb, nil)

	mock.ExpectQuery("SELECT count\\(\\*\\) FROM \"servers\" WHERE status = ?").
		WithArgs("On").
		WillReturnError(fmt.Errorf("db error"))

	num, err := repo.GetNumOnServers()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if num != 0 {
		t.Errorf("expected 0 ON servers on error, got %d", num)
	}
}

func TestGetNumOffServers_Success(t *testing.T) {
	gdb, mock, cleanup := repository.SetupMockDB(t)
	defer cleanup()

	repo := repository.NewServerInfoRepository(gdb, nil)

	mock.ExpectQuery("SELECT count\\(\\*\\) FROM \"servers\" WHERE status = ?").
		WithArgs("Off").
		WillReturnRows(mock.NewRows([]string{"count"}).AddRow(2))

	num, err := repo.GetNumOffServers()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if num != 2 {
		t.Errorf("expected 2 OFF servers, got %d", num)
	}
}

func TestGetNumOffServers_Error(t *testing.T) {
	gdb, mock, cleanup := repository.SetupMockDB(t)
	defer cleanup()

	repo := repository.NewServerInfoRepository(gdb, nil)

	mock.ExpectQuery("SELECT count\\(\\*\\) FROM \"servers\" WHERE status = ?").
		WithArgs("Off").
		WillReturnError(fmt.Errorf("db error"))

	num, err := repo.GetNumOffServers()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if num != 0 {
		t.Errorf("expected 0 OFF servers on error, got %d", num)
	}
}

func TestGetServerSumUpTimeRatio_Success(t *testing.T) {
	// Mock Elasticsearch client
	mockESC := new(MockESClient)

	// Prepare mock response body
	mockBuckets := []interface{}{
		map[string]interface{}{
			"doc_count": float64(4),
			"last_ping": map[string]interface{}{
				"hits": map[string]interface{}{
					"hits": []interface{}{
						map[string]interface{}{
							"_source": map[string]interface{}{
								"Status": "On",
							},
						},
					},
				},
			},
			"on_ping": map[string]interface{}{
				"total_on_ping_time": map[string]interface{}{
					"value": float64(10000),
				},
			},
			"off_ping": map[string]interface{}{
				"total_off_ping_time": map[string]interface{}{
					"value": float64(20000),
				},
			},
		},
	}
	mockAggs := map[string]interface{}{
		"id_bucket": map[string]interface{}{
			"buckets": mockBuckets,
		},
	}
	mockAnswer := map[string]interface{}{
		"aggregations": mockAggs,
	}

	respBody, _ := json.Marshal(mockAnswer)
	resp := &esapi.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader(respBody)),
	}

	mockESC.On("Search", mock.Anything, "ping_status", mock.Anything).Return(resp, nil)

	repo := repository.NewServerInfoRepository(nil, mockESC)

	start := time.Now().Add(-time.Hour).UTC().Format(time.RFC3339)
	end := time.Now().UTC().Format(time.RFC3339)

	ratio, err := repo.GetServerSumUpTimeRatio(start, end)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ratio == 0 {
		t.Errorf("expected non-zero ratio, got %v", ratio)
	}
}

func TestGetServerSumUpTimeRatio_InvalidStartTime(t *testing.T) {
	repo := repository.NewServerInfoRepository(nil, nil)
	_, err := repo.GetServerSumUpTimeRatio("invalid", time.Now().UTC().Format(time.RFC3339))
	if err == nil {
		t.Fatal("expected error for invalid start time, got nil")
	}
}

func TestGetServerSumUpTimeRatio_InvalidEndTime(t *testing.T) {
	repo := repository.NewServerInfoRepository(nil, nil)
	_, err := repo.GetServerSumUpTimeRatio(time.Now().UTC().Format(time.RFC3339), "invalid")
	if err == nil {
		t.Fatal("expected error for invalid end time, got nil")
	}
}

func TestGetServerSumUpTimeRatio_ESError(t *testing.T) {
	mockESC := new(MockESClient)
	mockESC.On("Search", mock.Anything, "ping_status", mock.Anything).Return(nil, assert.AnError)
	repo := repository.NewServerInfoRepository(nil, mockESC)

	start := time.Now().Add(-time.Hour).UTC().Format(time.RFC3339)
	end := time.Now().UTC().Format(time.RFC3339)

	_, err := repo.GetServerSumUpTimeRatio(start, end)
	if err == nil {
		t.Fatal("expected error for ES, got " + err.Error())
	}
}

func TestGetServerSumUpTimeRatio_ESReturnsErrorField(t *testing.T) {
	mockESC := new(MockESClient)
	mockAnswer := map[string]interface{}{
		"error":  "some error",
		"status": float64(500),
	}
	respBody, _ := json.Marshal(mockAnswer)
	resp := &esapi.Response{
		StatusCode: 500,
		Body:       io.NopCloser(bytes.NewReader(respBody)),
	}
	mockESC.On("Search", mock.Anything, "ping_status", mock.Anything).Return(resp, nil)
	repo := repository.NewServerInfoRepository(nil, mockESC)

	start := time.Now().Add(-time.Hour).UTC().Format(time.RFC3339)
	end := time.Now().UTC().Format(time.RFC3339)

	ratio, err := repo.GetServerSumUpTimeRatio(start, end)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ratio != 0 {
		t.Errorf("expected 0 ratio on ES error, got %v", ratio)
	}
}

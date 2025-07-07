package repository_test

import (
	"context"
	"fmt"
	"net/http"
	"server_administration_service/internal/repository"
	"testing"

	"github.com/elastic/go-elasticsearch/v9/esapi"
	"github.com/stretchr/testify/mock"
)

type MockESClient struct {
	mock.Mock
}

func (m *MockESClient) Do(ctx context.Context, req *http.Request) (*esapi.Response, error) {
	args := m.Called(ctx, req)
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
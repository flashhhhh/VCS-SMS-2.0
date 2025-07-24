package repository_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"server_administration_service/internal/repository"

	"github.com/elastic/go-elasticsearch/v9/esapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/DATA-DOG/go-sqlmock"
)

// Mock for ElasticsearchClient
type mockESC struct {
	mock.Mock
}

func (m *mockESC) Index(ctx context.Context, index string, data []byte) error {
	args := m.Called(ctx, index, data)
	return args.Error(0)
}

func (m *mockESC) Search(ctx context.Context, index string, buf bytes.Buffer) (*esapi.Response, error) {
	args := m.Called(ctx, index, buf)
	return args.Get(0).(*esapi.Response), args.Error(1)
}

func TestServerKafkaRepository_UpdateStatus_DBError(t *testing.T) {
	gdb, mock, cleanup := repository.SetupMockDB(t)
	defer cleanup()

	mockESC := new(mockESC)
	repo := repository.NewServerKafkaRepository(gdb, mockESC)

	serverID := "server-1"
	status := "active"

	// Simulate DB error
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"servers\" SET \"status\"").
		WithArgs(status, serverID).
		WillReturnError(errors.New("db error"))

	err := repo.UpdateStatus(serverID, status)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
}

func TestServerKafkaRepository_UpdateStatus_ESIndexError(t *testing.T) {
	gdb, mockDB, cleanup := repository.SetupMockDB(t)
	defer cleanup()

	mockESC := new(mockESC)
	repo := repository.NewServerKafkaRepository(gdb, mockESC)

	serverID := "server-1"
	status := "inactive"

	mockDB.ExpectBegin()
	mockDB.ExpectExec("UPDATE \"servers\" SET \"status\"").
		WithArgs(status, sqlmock.AnyArg(), serverID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockDB.ExpectRollback() // Expect rollback due to ES error

	mockESC.On("Index", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("es error"))

	err := repo.UpdateStatus(serverID, status)
	assert.Error(t, err)
}

func TestServerKafkaRepository_UpdateStatus_Success(t *testing.T) {
	gdb, mockDB, cleanup := repository.SetupMockDB(t)
	defer cleanup()

	mockESC := new(mockESC)
	repo := repository.NewServerKafkaRepository(gdb, mockESC)

	serverID := "server-2"
	status := "active"

	mockDB.ExpectBegin()
	mockDB.ExpectExec("UPDATE \"servers\" SET \"status\"").
		WithArgs(status, sqlmock.AnyArg(), serverID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockDB.ExpectCommit()

	mockESC.On("Index", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	err := repo.UpdateStatus(serverID, status)
	assert.NoError(t, err)
}
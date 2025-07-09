package repository_test

import (
	"errors"
	"regexp"
	"server_administration_service/internal/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetServerAddresses_Success(t *testing.T) {
	gdb, mock, cleanup := repository.SetupMockDB(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"server_id", "ipv4", "status"}).
		AddRow("srv1", "192.168.1.1", "active").
		AddRow("srv2", "192.168.1.2", "inactive")

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT "server_id","ipv4","status" FROM "servers"`)).
		WillReturnRows(rows)

	repo := repository.NewServerGRPCRepository(gdb)
	addresses, err := repo.GetServerAddresses()
	assert.NoError(t, err)
	assert.Len(t, addresses, 2)
	assert.Equal(t, "srv1", addresses[0].ServerID)
	assert.Equal(t, "192.168.1.2", addresses[1].IPv4)
}

func TestGetServerAddresses_Error(t *testing.T) {
	gdb, mock, cleanup := repository.SetupMockDB(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT server_id,ipv4,status FROM "servers"`)).
		WillReturnError(errors.New("db error"))

	repo := repository.NewServerGRPCRepository(gdb)
	addresses, err := repo.GetServerAddresses()
	assert.Error(t, err)
	assert.Nil(t, addresses)
}
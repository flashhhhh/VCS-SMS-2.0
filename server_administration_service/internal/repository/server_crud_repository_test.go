package repository_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"server_administration_service/internal/domain"
	"server_administration_service/internal/dto"
	"server_administration_service/internal/repository"
)

func TestCreateServer_Success(t *testing.T) {
	gdb, mock, cleanup := repository.SetupMockDB(t)
	defer cleanup()

	repo := repository.NewServerCRUDRepository(gdb)
	server := &domain.Server{
		ServerID:   "srv-1",
		ServerName: "TestServer",
		Status:     "On",
		IPv4:       "192.168.1.1",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "servers"`).
		WithArgs(
			server.ServerID,
			server.ServerName,
			server.Status,
			sqlmock.AnyArg(), // created_time
			sqlmock.AnyArg(), // last_updated
			server.IPv4,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	id, err := repo.CreateServer(server)
	assert.NoError(t, err)
	assert.Equal(t, server.ServerID, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateServer_FailDB(t *testing.T) {
	gdb, mock, cleanup := repository.SetupMockDB(t)
	defer cleanup()

	repo := repository.NewServerCRUDRepository(gdb)
	server := &domain.Server{
		ServerID: "server-1",
		ServerName: "Server 1",
		Status: "Off",
		IPv4: "192.168.1.1",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "servers"`).
		WithArgs(
			server.ServerID,
			server.ServerName,
			server.Status,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			server.IPv4,
		).
		WillReturnError(assert.AnError)
	mock.ExpectRollback()

	_, err := repo.CreateServer(server)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateServers_Success(t *testing.T) {
	gdb, mock, cleanup := repository.SetupMockDB(t)
	defer cleanup()

	repo := repository.NewServerCRUDRepository(gdb)
	servers := []domain.Server{
		{
			ServerID:   "srv-1",
			ServerName: "Server1",
			Status:     "On",
			IPv4:       "192.168.1.1",
		},
		{
			ServerID:   "srv-2",
			ServerName: "Server2",
			Status:     "Off",
			IPv4:       "192.168.1.2",
		},
	}

	// Build expected SQL
	expectedSQL := `INSERT INTO servers \(server_id, server_name, status, ipv4\) VALUES \('srv-1', 'Server1', 'On', '192.168.1.1'\), \('srv-2', 'Server2', 'Off', '192.168.1.2'\) ON CONFLICT DO NOTHING RETURNING \*`

	rows := sqlmock.NewRows([]string{"server_id", "server_name", "status", "ipv4"}).
		AddRow("srv-1", "Server1", "On", "192.168.1.1")

	mock.ExpectQuery(expectedSQL).WillReturnRows(rows)

	inserted, nonInserted, err := repo.CreateServers(servers)
	assert.NoError(t, err)
	assert.Len(t, inserted, 1)
	assert.Len(t, nonInserted, 1)
	assert.Equal(t, "srv-1", inserted[0].ServerID)
	assert.Equal(t, "srv-2", nonInserted[0].ServerID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateServers_FailDB(t *testing.T) {
	gdb, mock, cleanup := repository.SetupMockDB(t)
	defer cleanup()

	repo := repository.NewServerCRUDRepository(gdb)
	servers := []domain.Server{
		{
			ServerID:   "srv-1",
			ServerName: "Server1",
			Status:     "On",
			IPv4:       "192.168.1.1",
		},
	}

	expectedSQL := `INSERT INTO servers \(server_id, server_name, status, ipv4\) VALUES \('srv-1', 'Server1', 'On', '192.168.1.1'\) ON CONFLICT DO NOTHING RETURNING \*`
	mock.ExpectQuery(expectedSQL).WillReturnError(assert.AnError)

	inserted, nonInserted, err := repo.CreateServers(servers)
	assert.Error(t, err)
	assert.Nil(t, inserted)
	assert.Nil(t, nonInserted)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestViewServers_Success(t *testing.T) {
	gdb, mock, cleanup := repository.SetupMockDB(t)
	defer cleanup()

	repo := repository.NewServerCRUDRepository(gdb)

	// Prepare filter and expected query
	filter := &dto.ServerFilter{
		ServerID:   "srv-1",
		ServerName: "Test",
		Status:     "On",
		IPv4:       "192.168.1.1",
	}
	from := 0
	to := 10
	sortedColumn := "server_id"
	order := "asc"

	// Build expected SQL with LIKE and WHEREs
	mock.ExpectQuery(`SELECT \* FROM "servers" WHERE server_id = \$1 AND server_name LIKE \$2 AND status = \$3 AND ipv4 = \$4 ORDER BY server_id asc LIMIT \$5`).
		WithArgs(
			filter.ServerID,
			"%"+filter.ServerName+"%",
			filter.Status,
			filter.IPv4,
			to - from,
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"server_id", "server_name", "status", "ipv4"}).
				AddRow("srv-1", "TestServer", "On", "192.168.1.1"),
		)

	servers, err := repo.ViewServers(filter, from, to, sortedColumn, order)
	assert.NoError(t, err)
	assert.Len(t, servers, 1)
	assert.Equal(t, "srv-1", servers[0].ServerID)
	assert.Equal(t, "TestServer", servers[0].ServerName)
	assert.Equal(t, "On", servers[0].Status)
	assert.Equal(t, "192.168.1.1", servers[0].IPv4)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestViewServers_FailDB(t *testing.T) {
	gdb, mock, cleanup := repository.SetupMockDB(t)
	defer cleanup()

	repo := repository.NewServerCRUDRepository(gdb)

	// Prepare filter and expected query
	filter := &dto.ServerFilter{
		ServerID:   "srv-1",
		ServerName: "Test",
		Status:     "On",
		IPv4:       "192.168.1.1",
	}
	from := 0
	to := 10
	sortedColumn := "server_id"
	order := "asc"

	// Build expected SQL with LIKE and WHEREs
	mock.ExpectQuery(`SELECT \* FROM "servers" WHERE server_id = \$1 AND server_name LIKE \$2 AND status = \$3 AND ipv4 = \$4 ORDER BY server_id asc LIMIT \$5`).
		WithArgs(
			filter.ServerID,
			"%"+filter.ServerName+"%",
			filter.Status,
			filter.IPv4,
			to - from,
		).
		WillReturnError(assert.AnError)

	_, err := repo.ViewServers(filter, from, to, sortedColumn, order)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateServer_Success(t *testing.T) {
	gdb, mock, cleanup := repository.SetupMockDB(t)
	defer cleanup()

	repo := repository.NewServerCRUDRepository(gdb)
	serverID := "srv-1"
	updatedData := map[string]interface{}{
		"server_name": "UpdatedServer",
		"status":      "Off",
	}

	// Expect update
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "servers"`).
		WithArgs(
			updatedData["server_name"],
			updatedData["status"],
			sqlmock.AnyArg(), // last_updated
			serverID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.UpdateServer(serverID, updatedData)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateServer_FailDB(t *testing.T) {
	gdb, mock, cleanup := repository.SetupMockDB(t)
	defer cleanup()

	repo := repository.NewServerCRUDRepository(gdb)
	serverID := "srv-2"
	updatedData := map[string]interface{}{
		"server_name": "FailServer",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "servers"`).
		WithArgs(
			updatedData["server_name"],
			sqlmock.AnyArg(), // last_updated
			serverID,
		).
		WillReturnError(assert.AnError)
	mock.ExpectRollback()

	err := repo.UpdateServer(serverID, updatedData)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteServer_Success(t *testing.T) {
	gdb, mock, cleanup := repository.SetupMockDB(t)
	defer cleanup()

	repo := repository.NewServerCRUDRepository(gdb)
	serverID := "srv-1"

	// Expect update
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "servers"`).
		WithArgs(
			serverID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.DeleteServer(serverID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteServer_FailDB(t *testing.T) {
	gdb, mock, cleanup := repository.SetupMockDB(t)
	defer cleanup()

	repo := repository.NewServerCRUDRepository(gdb)
	serverID := "srv-1"

	// Expect update
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "servers"`).
		WithArgs(
			serverID,
		).
		WillReturnError(assert.AnError)
	mock.ExpectRollback()

	err := repo.DeleteServer(serverID)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
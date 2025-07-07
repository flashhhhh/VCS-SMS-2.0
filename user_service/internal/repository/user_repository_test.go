package repository_test

import (
	"errors"
	"testing"

	"user_service/internal/domain"
	"user_service/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn: db,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	cleanup := func() {
		db.Close()
	}
	return gormDB, mock, cleanup
}

func TestCreateUser_Success(t *testing.T) {
	db, mock, cleanup := SetupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRepository(db)

	user := &domain.User{
		ID:       "123",
		Username: "testuser",
		Password: "password",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users"`).
		WithArgs(user.ID, user.Username, user.Password).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(user.ID))
	mock.ExpectCommit()

	id, err := repo.CreateUser(user)
	assert.NoError(t, err)
	assert.Equal(t, "123", id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_Error(t *testing.T) {
	db, mock, cleanup := SetupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRepository(db)

	user := &domain.User{
		ID:       "123",
		Username: "testuser",
		Password: "password",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users"`).
		WithArgs(user.ID, user.Username, user.Password).
		WillReturnError(errors.New("insert error"))
	mock.ExpectRollback()

	id, err := repo.CreateUser(user)
	assert.Error(t, err)
	assert.Empty(t, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}


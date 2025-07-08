package repository_test

import (
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
		Name:	  "Test User",
		Email:	  "testuser@gmail.com",
		Role: 	  "user",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "users"`).
		WithArgs(user.ID, user.Username, user.Password, user.Name, user.Email, user.Role).
		WillReturnResult(sqlmock.NewResult(1, 1))
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
		Name:	  "Test User",
		Email:	  "testuser@gmail.com",
		Role: 	  "user",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "users"`).
		WithArgs(user.ID, user.Username, user.Password, user.Name, user.Email, user.Role).
		WillReturnError(assert.AnError)
	mock.ExpectRollback()

	id, err := repo.CreateUser(user)
	assert.Error(t, err)
	assert.Empty(t, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLogin_Success(t *testing.T) {
	db, mock, cleanup := SetupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRepository(db)

	expectedUser := &domain.User{
		ID:       "123",
		Username: "testuser",
		Password: "password",
		Name:	  "Test User",
		Email:	  "testuser@gmail.com",
		Role: 	  "user",
	}

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE username = \$1 ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(expectedUser.Username, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "name", "email", "role"}).
			AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Password, expectedUser.Name, expectedUser.Email, expectedUser.Role))

	user, err := repo.Login("testuser")
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Username, user.Username)
	assert.Equal(t, expectedUser.Password, user.Password)
	assert.Equal(t, expectedUser.Name, user.Name)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Role, user.Role)
	
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLogin_FailDB(t *testing.T) {
	db, mock, cleanup := SetupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRepository(db)

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE username = \$1 ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs("testuser", 1).
		WillReturnError(assert.AnError)

	_, err := repo.Login("testuser")
	assert.Error(t, err)
	
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByID_Success(t *testing.T) {
	db, mock, cleanup := SetupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRepository(db)

	expectedUser := &domain.User{
		ID:       "123",
		Username: "testuser",
		Password: "password",
		Name:     "Test User",
		Email:    "testuser@gmail.com",
		Role:     "user",
	}

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE id = \$1 ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(expectedUser.ID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "name", "email", "role"}).
			AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Password, expectedUser.Name, expectedUser.Email, expectedUser.Role))

	user, err := repo.GetUserByID("123")
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Username, user.Username)
	assert.Equal(t, expectedUser.Password, user.Password)
	assert.Equal(t, expectedUser.Name, user.Name)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Role, user.Role)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByID_NotFound(t *testing.T) {
	db, mock, cleanup := SetupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRepository(db)

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE id = \$1 ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs("notfound", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "name", "email", "role"})) // no rows

	user, err := repo.GetUserByID("notfound")
	assert.Error(t, err)
	assert.Nil(t, user)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByID_DBError(t *testing.T) {
	db, mock, cleanup := SetupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRepository(db)

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE id = \$1 ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs("123", 1).
		WillReturnError(assert.AnError)

	user, err := repo.GetUserByID("123")
	assert.Error(t, err)
	assert.Nil(t, user)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllUsers_Success(t *testing.T) {
	db, mock, cleanup := SetupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRepository(db)

	expectedUsers := []*domain.User{
		{
			ID:       "1",
			Username: "user1",
			Password: "pass1",
			Name:     "User One",
			Email:    "user1@example.com",
			Role:     "admin",
		},
		{
			ID:       "2",
			Username: "user2",
			Password: "pass2",
			Name:     "User Two",
			Email:    "user2@example.com",
			Role:     "user",
		},
	}

	rows := sqlmock.NewRows([]string{"id", "username", "password", "name", "email", "role"}).
		AddRow(expectedUsers[0].ID, expectedUsers[0].Username, expectedUsers[0].Password, expectedUsers[0].Name, expectedUsers[0].Email, expectedUsers[0].Role).
		AddRow(expectedUsers[1].ID, expectedUsers[1].Username, expectedUsers[1].Password, expectedUsers[1].Name, expectedUsers[1].Email, expectedUsers[1].Role)

	mock.ExpectQuery(`SELECT \* FROM "users"`).
		WillReturnRows(rows)

	users, err := repo.GetAllUsers()
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, expectedUsers[0].ID, users[0].ID)
	assert.Equal(t, expectedUsers[1].ID, users[1].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllUsers_Empty(t *testing.T) {
	db, mock, cleanup := SetupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRepository(db)

	rows := sqlmock.NewRows([]string{"id", "username", "password", "name", "email", "role"})
	mock.ExpectQuery(`SELECT \* FROM "users"`).
		WillReturnRows(rows)

	users, err := repo.GetAllUsers()
	assert.NoError(t, err)
	assert.Len(t, users, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllUsers_DBError(t *testing.T) {
	db, mock, cleanup := SetupMockDB(t)
	defer cleanup()

	repo := repository.NewUserRepository(db)

	mock.ExpectQuery(`SELECT \* FROM "users"`).
		WillReturnError(assert.AnError)

	users, err := repo.GetAllUsers()
	assert.Error(t, err)
	assert.Nil(t, users)
	assert.NoError(t, mock.ExpectationsWereMet())
}
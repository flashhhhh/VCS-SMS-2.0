package service_test

import (
	"testing"
	"user_service/internal/domain"
	"user_service/internal/service"

	"github.com/flashhhhh/pkg/hash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository implements repository.UserRepository for testing
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *domain.User) (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}
func (m *MockUserRepository) Login(username string) (*domain.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}
func (m *MockUserRepository) GetUserByID(id string) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}
func (m *MockUserRepository) GetAllUsers() ([]*domain.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func TestUserService_CreateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userSvc := service.NewUserService(mockRepo)

	username := "testuser"
	password := "password123"
	name := "Test User"
	email := "test@example.com"
	role := "admin"

	mockRepo.On("CreateUser", mock.Anything).Return("123", nil)

	userID, err := userSvc.CreateUser(username, password, name, email, role)

	assert.NoError(t, err)
	assert.Equal(t, "123", userID)
}

func TestUserService_CreateUser_FailDB(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userSvc := service.NewUserService(mockRepo)

	username := "testuser"
	password := "password123"
	name := "Test User"
	email := "test@example.com"
	role := "admin"

	mockRepo.On("CreateUser", mock.Anything).Return("", assert.AnError)

	userID, err := userSvc.CreateUser(username, password, name, email, role)

	assert.Error(t, err)
	assert.Equal(t, "", userID)
}

func TestUserService_Login_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userSvc := service.NewUserService(mockRepo)

	username := "testuser"
	password := "password123"
	hashedPassword := hash.HashString(password)
	user := &domain.User{
		ID:       "1",
		Username: username,
		Password: hashedPassword,
		Name:     "Test User",
		Email:    "test@example.com",
		Role:     "admin",
	}

	mockRepo.On("Login", username).Return(user, nil)

	token, err := userSvc.Login(username, password)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userSvc := service.NewUserService(mockRepo)

	username := "testuser"
	password := "password123"

	mockRepo.On("Login", username).Return(nil, assert.AnError)

	token, err := userSvc.Login(username, password)
	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestUserService_Login_InvalidPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userSvc := service.NewUserService(mockRepo)

	username := "testuser"
	password := "wrongpassword"
	hashedPassword := "hashedpassword"
	user := &domain.User{
		ID:       "1",
		Username: username,
		Password: hashedPassword,
		Name:     "Test User",
		Email:    "test@example.com",
		Role:     "admin",
	}

	mockRepo.On("Login", username).Return(user, nil)

	token, err := userSvc.Login(username, password)
	assert.Error(t, err)
	assert.Equal(t, "invalid password", err.Error())
	assert.Empty(t, token)
}

func TestUserService_GetUserByID_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userSvc := service.NewUserService(mockRepo)

	userID := "1"
	user := &domain.User{
		ID:       userID,
		Username: "testuser",
		Password: "hashedpassword",
		Name:     "Test User",
		Email:    "test@example.com",
		Role:     "admin",
	}

	mockRepo.On("GetUserByID", userID).Return(user, nil)

	result, err := userSvc.GetUserByID(userID)
	assert.NoError(t, err)
	assert.Equal(t, user, result)
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userSvc := service.NewUserService(mockRepo)

	userID := "2"
	mockRepo.On("GetUserByID", userID).Return(nil, assert.AnError)

	result, err := userSvc.GetUserByID(userID)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "User "+userID+" not found")
}

func TestUserService_GetAllUsers_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userSvc := service.NewUserService(mockRepo)

	users := []*domain.User{
		{
			ID:       "1",
			Username: "user1",
			Password: "hashed1",
			Name:     "User One",
			Email:    "one@example.com",
			Role:     "admin",
		},
		{
			ID:       "2",
			Username: "user2",
			Password: "hashed2",
			Name:     "User Two",
			Email:    "two@example.com",
			Role:     "user",
		},
	}

	mockRepo.On("GetAllUsers").Return(users, nil)

	result, err := userSvc.GetAllUsers()
	assert.NoError(t, err)
	assert.Equal(t, users, result)
}

func TestUserService_GetAllUsers_Fail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userSvc := service.NewUserService(mockRepo)

	mockRepo.On("GetAllUsers").Return(nil, assert.AnError)

	result, err := userSvc.GetAllUsers()
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Failed to retrieve all users")
}
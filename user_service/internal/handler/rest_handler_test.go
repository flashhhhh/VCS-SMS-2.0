package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"user_service/internal/domain"
	"user_service/internal/handler"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService implements service.UserService for testing
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(username, password, name, email, role string) (string, error) {
	args := m.Called(username, password, name, email, role)
	return args.String(0), args.Error(1)
}
func (m *MockUserService) Login(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}
func (m *MockUserService) GetUserByID(userID string) (*domain.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}
func (m *MockUserService) GetAllUsers() ([]*domain.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func TestRestHandler_CreateUser_Success(t *testing.T) {
	mockService := new(MockUserService)
	h := handler.NewRestHandler(mockService)

	body := map[string]interface{}{
		"username": "testuser",
		"password": "testpass",
		"name":     "Test User",
		"email":    "test@example.com",
		"role":     "admin",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()

	mockService.On("CreateUser", "testuser", "testpass", "Test User", "test@example.com", "admin").Return("123", nil)

	h.CreateUser(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var respBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respBody)
	assert.Equal(t, "User created successfully", respBody["message"])
	assert.Equal(t, "123", respBody["userID"])
}

func TestRestHandler_CreateUser_BadRequest(t *testing.T) {
	mockService := new(MockUserService)
	h := handler.NewRestHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	h.CreateUser(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestRestHandler_CreateUser_ServiceError(t *testing.T) {
	mockService := new(MockUserService)
	h := handler.NewRestHandler(mockService)

	body := map[string]interface{}{
		"username": "testuser",
		"password": "testpass",
		"name":     "Test User",
		"email":    "test@example.com",
		"role":     "admin",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()

	mockService.On("CreateUser", "testuser", "testpass", "Test User", "test@example.com", "admin").Return("", errors.New("db error"))

	h.CreateUser(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestRestHandler_Login_Success(t *testing.T) {
	mockService := new(MockUserService)
	h := handler.NewRestHandler(mockService)

	body := map[string]interface{}{
		"username": "testuser",
		"password": "testpass",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()

	mockService.On("Login", "testuser", "testpass").Return("token123", nil)

	h.Login(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var respBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respBody)
	assert.Equal(t, "User logged in successfully", respBody["message"])
	assert.Equal(t, "token123", respBody["token"])
}

func TestRestHandler_Login_BadRequest(t *testing.T) {
	mockService := new(MockUserService)
	h := handler.NewRestHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte("bad json")))
	w := httptest.NewRecorder()

	h.Login(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestRestHandler_Login_InvalidPassword(t *testing.T) {
	mockService := new(MockUserService)
	h := handler.NewRestHandler(mockService)

	body := map[string]interface{}{
		"username": "testuser",
		"password": "wrongpass",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()

	mockService.On("Login", "testuser", "wrongpass").Return("", errors.New("Invalid password"))

	h.Login(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestRestHandler_Login_ServiceError(t *testing.T) {
	mockService := new(MockUserService)
	h := handler.NewRestHandler(mockService)

	body := map[string]interface{}{
		"username": "testuser",
		"password": "testpass",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()

	mockService.On("Login", "testuser", "testpass").Return("", errors.New("db error"))

	h.Login(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestRestHandler_GetUserByID_Success(t *testing.T) {
	mockService := new(MockUserService)
	h := handler.NewRestHandler(mockService)

	user := &domain.User{ID: "123", Username: "testuser"}
	mockService.On("GetUserByID", "123").Return(user, nil)

	req := httptest.NewRequest(http.MethodGet, "/user?userID=123", nil)
	req.Header.Set("userRole", "admin")
	w := httptest.NewRecorder()

	h.GetUserByID(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var respBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respBody)
	assert.Equal(t, "User retrieved successfully", respBody["message"])
}

func TestRestHandler_GetUserByID_Forbidden(t *testing.T) {
	mockService := new(MockUserService)
	h := handler.NewRestHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/user?userID=456", nil)
	req.Header.Set("userID", "123")
	w := httptest.NewRecorder()

	h.GetUserByID(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestRestHandler_GetUserByID_NotFound(t *testing.T) {
	mockService := new(MockUserService)
	h := handler.NewRestHandler(mockService)

	mockService.On("GetUserByID", "999").Return(nil, errors.New("User 999 not found"))

	req := httptest.NewRequest(http.MethodGet, "/user?userID=999", nil)
	req.Header.Set("userRole", "admin")
	w := httptest.NewRecorder()

	h.GetUserByID(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestRestHandler_GetUserByID_ServiceError(t *testing.T) {
	mockService := new(MockUserService)
	h := handler.NewRestHandler(mockService)

	mockService.On("GetUserByID", "888").Return(nil, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/user?userID=888", nil)
	req.Header.Set("userRole", "admin")
	w := httptest.NewRecorder()

	h.GetUserByID(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestRestHandler_GetAllUsers_Success(t *testing.T) {
	mockService := new(MockUserService)
	h := handler.NewRestHandler(mockService)

	users := []*domain.User{
		{ID: "1", Username: "user1"},
		{ID: "2", Username: "user2"},
	}
	mockService.On("GetAllUsers").Return(users, nil)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	h.GetAllUsers(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var respBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respBody)
	assert.Equal(t, "All users retrieved successfully", respBody["message"])
}

func TestRestHandler_GetAllUsers_ServiceError(t *testing.T) {
	mockService := new(MockUserService)
	h := handler.NewRestHandler(mockService)

	mockService.On("GetAllUsers").Return(nil, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	h.GetAllUsers(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
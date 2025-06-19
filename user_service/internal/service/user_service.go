package service

import (
	"errors"
	"time"
	"user_service/internal/domain"
	"user_service/internal/repository"

	"github.com/flashhhhh/pkg/hash"
	"github.com/flashhhhh/pkg/jwt"
	"github.com/flashhhhh/pkg/logging"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(username, password, name, email, role string) (string, error)
	Login(username, password string) (string, error)
	GetUserByID(id string) (*domain.User, error)
	GetAllUsers() ([]*domain.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	logging.LogMessage("user_service", "Initializing UserService", "INFO")

	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) CreateUser(username, password, name, email, role string) (string, error) {
	hashedPassword := hash.HashString(password)

	user := &domain.User{
		ID : uuid.New().String(),
		Username: username,
		Password: hashedPassword,
		Name:     name,
		Email:    email,
		Role:     role,
	}

	userID, err := s.userRepository.CreateUser(user)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (s *userService) Login(username, password string) (string, error) {
	user, err := s.userRepository.Login(username)
	if err != nil {
		return "", err
	}

	if !hash.CompareHashAndString(user.Password, password) {
		return "", errors.New("invalid password")
	}

	token, _ := jwt.GenerateToken(
		map[string]any{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		}, time.Hour)
	
	return token, nil
}

func (s *userService) GetUserByID(id string) (*domain.User, error) {
	user, err := s.userRepository.GetUserByID(id)
	if err != nil {
		return nil, errors.New("User " + id + " not found: " + err.Error())
	}
	return user, nil
}

func (s *userService) GetAllUsers() ([]*domain.User, error) {
	users, err := s.userRepository.GetAllUsers()
	if err != nil {
		return nil, errors.New("Failed to retrieve all users: " + err.Error())
	}
	return users, nil
}
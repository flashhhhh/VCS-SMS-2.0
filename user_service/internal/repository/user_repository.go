package repository

import (
	"user_service/internal/domain"

	"github.com/flashhhhh/pkg/logging"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *domain.User) (string, error)
	Login(username string) (*domain.User, error)
	GetUserByID(id string) (*domain.User, error)
	GetAllUsers() ([]*domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	logging.LogMessage("user_service", "Initializing UserRepository", "INFO")

	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(user *domain.User) (string, error) {
	err := r.db.Create(user).Error
	if err != nil {
		return "", err
	}
	
	return user.ID, nil
}

func (r *userRepository) Login(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByID(id string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetAllUsers() ([]*domain.User, error) {
	var users []*domain.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
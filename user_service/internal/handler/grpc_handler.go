package handler

import (
	"context"
	"errors"
	"user_service/internal/service"
	"user_service/pb"

	"github.com/flashhhhh/pkg/logging"
)

type GrpcHandler struct {
	userService service.UserService
	pb.UnimplementedUserServiceServer
}

func NewGrpcHandler(userService service.UserService) *GrpcHandler {
	return &GrpcHandler{
		userService: userService,
	}
}

func (grpcHandler *GrpcHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	logging.LogMessage("user_service", "Creating new user", "INFO")
	logging.LogMessage("user_service", "username: "+req.Username+", password: "+req.Password+", name: "+req.Name+", email: "+req.Email+", role: "+req.Role, "DEBUG")

	userID, err := grpcHandler.userService.CreateUser(req.Username, req.Password, req.Name, req.Email, req.Role)
	if err != nil {
		if err.Error() == "user already exists" {
			logging.LogMessage("user_service", "User already exists: "+err.Error(), "ERROR")
			return nil, errors.New("user already exists")
		} else {
			logging.LogMessage("user_service", "Failed to create user: "+err.Error(), "ERROR")
			return nil, err
		}
	}

	logging.LogMessage("user_service", "User created successfully with ID: "+userID, "INFO")
	return &pb.CreateUserResponse{
		UserID:  userID,
	}, nil
}

func (grpcHandler *GrpcHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	logging.LogMessage("user_service", "User login attempt", "INFO")
	logging.LogMessage("user_service", "username: "+req.Username+", password: "+req.Password, "DEBUG")

	token, err := grpcHandler.userService.Login(req.Username, req.Password)
	if err != nil {
		if err.Error() == "invalid password" {
			logging.LogMessage("user_service", "Invalid password: "+err.Error(), "ERROR")
			return nil, errors.New("invalid password")
		}

		logging.LogMessage("user_service", "Failed to login: "+err.Error(), "ERROR")
		return nil, err
	}

	logging.LogMessage("user_service", "User logged in successfully", "INFO")
	logging.LogMessage("user_service", "Generated token: "+token, "DEBUG")

	return &pb.LoginResponse{
		Token: token,
	}, nil
}

func (grpcHandler *GrpcHandler) GetUserByID(ctx context.Context, req *pb.IDRequest) (*pb.UserResponse, error) {
	logging.LogMessage("user_service", "Fetching user by ID", "INFO")
	logging.LogMessage("user_service", "userID: "+req.Id, "DEBUG")

	user, err := grpcHandler.userService.GetUserByID(req.Id)
	if err != nil {
		logging.LogMessage("user_service", "Failed to fetch user: "+err.Error(), "ERROR")
		return nil, err
	}

	logging.LogMessage("user_service", "User fetched successfully", "INFO")
	return &pb.UserResponse{
		UserID: user.ID,
		Username: user.Username,
		Name:     user.Name,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}

func (grpcHandler *GrpcHandler) GetAllUsers(ctx context.Context, req *pb.EmptyRequest) (*pb.UsersResponse, error) {
	logging.LogMessage("user_service", "Fetching all users", "INFO")

	users, err := grpcHandler.userService.GetAllUsers()
	if err != nil {
		logging.LogMessage("user_service", "Failed to fetch users: "+err.Error(), "ERROR")
		return nil, err
	}

	logging.LogMessage("user_service", "All users fetched successfully", "INFO")

	logging.LogMessage("user_service", "Converting users to protobuf form", "DEBUG")
	var pbUsers []*pb.UserResponse
	for _, user := range users {
		pbUsers = append(pbUsers, &pb.UserResponse{
			UserID:   user.ID,
			Username: user.Username,
			Name:     user.Name,
			Email:    user.Email,
			Role:     user.Role,
		})
	}

	logging.LogMessage("user_service", "Successfully converting users to protobuf form", "DEBUG")

	return &pb.UsersResponse{
		Users: pbUsers,
	}, nil
}
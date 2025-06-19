package main

import (
	"context"
	"log"
	"user_service/pb"

	"github.com/flashhhhh/pkg/env"
	"google.golang.org/grpc"
)

func main() {
	// Connect to the gRPC server
	grpcServerAddress := env.GetEnv("USER_SERVER_HOST", "localhost") + ":" + env.GetEnv("GRPC_USER_SERVER_PORT", "50051")

	conn, err := grpc.Dial(grpcServerAddress, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	client := pb.NewUserServiceClient(conn)

	// Example usage of the client
	
	/*
		Create a user
	*/
	// _, err = client.CreateUser(context.Background(), &pb.CreateUserRequest{
	// 	Username: "testuser2",
	// 	Password: "password",
	// 	Name:     "Test User 2",
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println("User created successfully!")

	/*
		Log in a user
	*/
	token, err := client.Login(context.Background(), &pb.LoginRequest{
		Username: "johndoe",
		Password: "password123",
	})
	if err != nil {
		panic(err)
	}
	log.Println("Login successful! Token:", token.Token)

	/*
		Get user by ID
	*/
	user, err := client.GetUserByID(context.Background(), &pb.IDRequest{
		Id: "b825a743-6b2a-484f-8fc4-25915c468a96",
	})
	if err != nil {
		panic(err)
	}
	log.Println("User details:", user)

	/*
		Get all users
	*/
	users, err := client.GetAllUsers(context.Background(), &pb.EmptyRequest{})
	if err != nil {
		panic(err)
	}
	log.Println("All users:")
	for _, user := range users.Users {
		log.Printf("ID: %s, Username: %s, Name: %s, Email: %s, Role: %s\n", user.UserID, user.Username, user.Name, user.Email, user.Role)
	}
}
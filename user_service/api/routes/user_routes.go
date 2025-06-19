package api

import (
	"net/http"
	"user_service/api/middlewares"
	"user_service/internal/handler"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, userHandler handler.RestHandler) {
	r.Handle("/create", middlewares.AdminMiddleware(http.HandlerFunc(userHandler.CreateUser))).Methods("POST")
	r.HandleFunc("/login", userHandler.Login).Methods("POST")
	r.Handle("/getUserByID", middlewares.UserMiddleware(http.HandlerFunc(userHandler.GetUserByID))).Methods("GET")
	r.Handle("/getAllUsers", middlewares.AdminMiddleware(http.HandlerFunc(userHandler.GetAllUsers))).Methods("GET")
}
package routes

import (
	"net/http"
	// "server_administration_service/api/middlewares"
	"server_administration_service/internal/handler"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, serverHandler handler.ServerRestHandler) {
	r.Handle("/create", http.HandlerFunc(serverHandler.CreateServer)).Methods("POST")
	r.Handle("/view", http.HandlerFunc(serverHandler.ViewServers)).Methods("GET")
	r.Handle("/update", http.HandlerFunc(serverHandler.UpdateServer)).Methods("PUT")
	r.Handle("/delete", http.HandlerFunc(serverHandler.DeleteServer)).Methods("DELETE")
	r.Handle("/import", http.HandlerFunc(serverHandler.ImportServers)).Methods("POST")
	r.Handle("/export", http.HandlerFunc(serverHandler.ExportServers)).Methods("GET")
}
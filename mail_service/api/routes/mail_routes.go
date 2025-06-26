package routes

import (
	"mail_service/api/middlewares"
	"mail_service/internal/handler"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, mailHandler handler.MailHandler) {
	r.Handle("/send", middlewares.AdminMiddleware(http.HandlerFunc(mailHandler.SendEmail))).Methods("POST")
}
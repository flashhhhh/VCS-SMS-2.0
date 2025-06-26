package handler

import (
	"mail_service/internal/service"
	"net/http"

	"github.com/flashhhhh/pkg/logging"
)

type MailHandler interface {
	SendEmail(w http.ResponseWriter, r *http.Request)
}

type mailHandler struct {
	mailService service.MailService
}

func NewMailHandler(mailService service.MailService) MailHandler {
	return &mailHandler{
		mailService: mailService,
	}
}

func (h *mailHandler) SendEmail(w http.ResponseWriter, r *http.Request) {
	to := r.URL.Query().Get("to")
	startTime := r.URL.Query().Get("start_time")
	endTime := r.URL.Query().Get("end_time")

	if to == "" {
		logging.LogMessage("mail_service", "Missing 'to' parameter", "ERROR")
		http.Error(w, "Missing 'to' parameter", http.StatusBadRequest)
		return
	}

	if startTime == "" {
		logging.LogMessage("mail_service", "Missing 'start_time' parameter", "ERROR")
		http.Error(w, "Missing 'start_time' parameter", http.StatusBadRequest)
		return
	}

	if endTime == "" {
		logging.LogMessage("mail_service", "Missing 'end_time' parameter", "ERROR")
		http.Error(w, "Missing 'end_time' parameter", http.StatusBadRequest)
		return
	}

	err := h.mailService.SendServersReportEmail(to, startTime, endTime)
	if err != nil {
		logging.LogMessage("mail_service", "Failed to send emails: "+err.Error(), "ERROR")
		http.Error(w, "Failed to send emails: "+err.Error(), http.StatusInternalServerError)
		return
	}

	logging.LogMessage("mail_service", "Email sent successfully!", "INFO")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Emails sent successfully"))
}
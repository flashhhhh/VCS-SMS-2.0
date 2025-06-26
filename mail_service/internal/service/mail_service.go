package service

import (
	"fmt"
	mailsending "mail_service/infrastructure/mail_sending"
	"mail_service/internal/repository"
	"time"

	"github.com/flashhhhh/pkg/logging"
)

type MailService interface {
	SendServersReportEmail(to string, startTime, endTime string) (error)
}

type mailService struct {
	mailSending mailsending.MailSending
	mailGRPCClientRepository repository.MailGRPCClientRepository
}

func NewMailService(mailSending mailsending.MailSending, mailGRPCClientRepository repository.MailGRPCClientRepository) MailService {
	return &mailService{
		mailSending: mailSending,
		mailGRPCClientRepository: mailGRPCClientRepository,
	}
}

func (ms *mailService) SendServersReportEmail(to string, startTime, endTime string) (error) {
	numServers, numOnServers, numOffServers, meanUpTimeRatio, err := ms.mailGRPCClientRepository.GetServersInformation(startTime, endTime)
	if err != nil {
		logging.LogMessage("mail_service", "Failed to get server information in range [" + startTime + ", " + endTime + "]. Err: " + err.Error(), "ERROR")
		return err
	}

	subject := "Daily Server Status Report for " + time.Now().Format("2006-01-02")
	body := fmt.Sprintf("Dear server administrator,\n\nThe server status is as follows:\n\nTotal servers: %d\nServers on: %d\nServers off: %d\nMean uptime rate: %.2f%%\n\nBest regards,\nYour Server Monitoring System", numServers, numOnServers, numOffServers, meanUpTimeRatio)

	err = ms.mailSending.SendEmail(to, subject, body)
	if err != nil {
		logging.LogMessage("mail_service", "Cannot send email for days in range [" + startTime + ", " + endTime + "]. Err: " + err.Error(), "ERROR")
		return err
	}

	logging.LogMessage("mail_service", "Send email successfully for days in range [" + startTime + ", " + endTime + "]!", "INFO")
	return nil
}
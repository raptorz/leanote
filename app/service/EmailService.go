package service

import (
	"github.com/leanote/leanote/app/info"
)

type EmailService struct {
}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (this *EmailService) SendEmail(to, subject, body string) (ok bool, e string) {
	return true, ""
}

func (this *EmailService) RegisterSendActiveEmail(userInfo info.User, email string) bool {
	return true
}

func (this *EmailService) UpdateEmailSendActiveEmail(userInfo info.User, email string) (ok bool, msg string) {
	return true, ""
}

func (this *EmailService) FindPwdSendEmail(token, email string) (ok bool, msg string) {
	return true, ""
}

func (this *EmailService) SendInviteEmail(userInfo info.User, email, content string) bool {
	return true
}

func (this *EmailService) SendCommentEmail(note info.Note, comment info.BlogComment, userId, content string) bool {
	return true
}

func (this *EmailService) SendEmailToUsers(users []info.User, subject, body string) (ok bool, msg string) {
	return true, ""
}

func (this *EmailService) SendEmailToEmails(emails []string, subject, body string) (ok bool, msg string) {
	return true, ""
}

func (this *EmailService) AddEmailLog(email, subject, body string, ok bool, msg string) {
}

func (this *EmailService) DeleteEmails(ids []string) bool {
	return true
}

func (this *EmailService) ListEmailLogs(pageNumber, pageSize int, sortField string, isAsc bool, email string) (page info.Page, emailLogs []info.EmailLog) {
	return info.Page{}, []info.EmailLog{}
}

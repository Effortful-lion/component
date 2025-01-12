package email

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type EmailSender struct {
    SMTPHost  string
    SMTPPort  int
    Username  string
    Password  string
    FromEmail string
}

// 初始化邮件发送器
func NewEmailSender(smtpHost string, smtpPort int, username, password, fromEmail string) *EmailSender {
    return &EmailSender{
        SMTPHost:  smtpHost,
        SMTPPort:  smtpPort,
        Username:  username,
        Password:  password,
        FromEmail: fromEmail,
    }
}

// 发送邮件
func (es *EmailSender) SendEmail(toEmail, subject, body string) error {
    m := gomail.NewMessage()
    m.SetHeader("From", es.FromEmail)
    m.SetHeader("To", toEmail)
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", body)

    d := gomail.NewDialer(es.SMTPHost, es.SMTPPort, es.Username, es.Password)
    err := d.DialAndSend(m)
    if err != nil {
        fmt.Println("Error sending email:", err) // 添加错误日志
    }
    return err
}
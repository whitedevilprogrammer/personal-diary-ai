package services

import (
	"fmt"
	"net/smtp"
	"os"
)

type EmailSender struct {
	smtpHost string
	smtpPort string
	username string
	password string
	from     string
}

func NewEmailSender() *EmailSender {
	return &EmailSender{
		smtpHost: os.Getenv("SMTP_HOST"),
		smtpPort: os.Getenv("SMTP_PORT"),
		username: os.Getenv("SMTP_USERNAME"),
		password: os.Getenv("SMTP_PASSWORD"),
		from:     os.Getenv("EMAIL_FROM"),
	}
}

func (e *EmailSender) Send(to, subject, plainBody, htmlBody string) error {
	auth := smtp.PlainAuth("", e.username, e.password, e.smtpHost)

	headers := map[string]string{
		"From":         e.from,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": `multipart/alternative; boundary="mixed-boundary"`,
	}

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n--mixed-boundary\r\n"
	message += "Content-Type: text/plain; charset=UTF-8\r\n\r\n" + plainBody
	message += "\r\n--mixed-boundary\r\n"
	message += "Content-Type: text/html; charset=UTF-8\r\n\r\n" + htmlBody
	message += "\r\n--mixed-boundary--"

	addr := fmt.Sprintf("%s:%s", e.smtpHost, e.smtpPort)
	return smtp.SendMail(addr, auth, e.from, []string{to}, []byte(message))
}

func (e *EmailSender) SendPasswordResetEmail(to string, resetToken string) error {
	resetURL := os.Getenv("RESET_URL_BASE") + "?token=" + resetToken

	subject := "🔒 Reset your password"

	plainText := fmt.Sprintf("Click the link to reset your password: %s", resetURL)

	htmlBody := fmt.Sprintf(`
		<html>
		<body>
			<h2>Reset Password</h2>
			<p>Click the button below to reset your password:</p>
			<a href="%s" style="display:inline-block; padding:10px 20px; background:#007BFF; color:white; text-decoration:none; border-radius:5px;">Reset Password</a>
			<p>If the button doesn't work, copy and paste this link into your browser:</p>
			<p><a href="%s">%s</a></p>
		</body>
		</html>`, resetURL, resetURL, resetURL)

	return e.Send(to, subject, plainText, htmlBody)
}

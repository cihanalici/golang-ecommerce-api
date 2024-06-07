package util

//implement  .SendPasswordResetMail(user.Email, resetPasswordLink)

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

// Mailer is the interface that wraps the basic Send method
func SendResetEmail(email, token string) error {
	fmt.Println("token: ", token)
	m := gomail.NewMessage()
	m.SetHeader("From", "no-reply@myapp.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Password Reset")
	m.SetBody("text/html", fmt.Sprintf("To reset your password, please click the following link: <a href=\"http://yourapp.com/reset-password?token=%s\">Reset Password</a>", token))

	// MailHog kullanarak e-posta göndermek için güncellendi
	d := gomail.NewDialer("localhost", 1025, "", "")

	return d.DialAndSend(m)
}

package mail

import (
	"fmt"
	"net/smtp"

	"github.com/night1010/everhealth/config"
	"github.com/jordan-wright/email"
)

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
	smtpTestServer    = "localhost:1025"
)

type SmtpGmail interface {
	SendEmail(string, string, bool) error
	SendEmailTest(string, string, bool) error
}

type smtpGmail struct {
	name       string
	address    string
	password   string
	prefixLink string
}

func NewSmtpGmail() SmtpGmail {
	emailConfig := config.NewEmailConfig()
	return &smtpGmail{
		name:       emailConfig.Name,
		address:    emailConfig.Address,
		password:   emailConfig.Password,
		prefixLink: emailConfig.PrefixLink,
	}
}

func emailVerifyContent(link string) (subject, content string) {
	subject = "Account Verification"
	content = fmt.Sprintf(`
	<div style="background-color: #F2F2F2; padding: 5px; border-radius: 0.5rem; display: grid; grid-template-columns: 1fr; gap: 2rem; align-items: center; justify-items: center; height: 100vh;">
		<div style="display: grid; grid-template-columns: 1fr; background-color: white; width: 100%%; border-radius: 0.5rem; align-items: center; gap: 2rem; padding: 2rem; margin: auto; text-align: center;">
			<div>
				<img style="height: auto; width: 10rem; object-fit: contain; margin-top: 2rem;" src="https://everhealth-asset.irfancen.com/assets/eh.png" alt="Everhealth logo" />
			</div>

			<div style="width: 25rem; margin: auto;">
				<h1 style="color: black; margin: 0;">Hello!</h1>
				<p>Thank you for joining <span style="color: #36A5B2; font-weight: bold;">Everhealth!</span></p>
				<br />
				<p>Please click the link to start your verification process!</p>
				<br />
				<a href="%s" style="text-decoration: none; color: white; background-color: #36A5B2; padding: 10px 20px; border-radius: 0.3rem; font-weight: bold; display: inline-block;">Verification Process</a>
				<br />
				<p>Best,</p>
				<p style="margin-bottom: 3rem; color: #36A5B2; font-weight: bold;">Everhealth</p>
			</div>
		</div>
	</div>
	`, link)
	return subject, content
}

func emailForgotPasswordContent(link string) (subject, content string) {
	subject = "Change Password Verification"
	content = fmt.Sprintf(`
	<div style="background-color: #F2F2F2; padding: 5px; border-radius: 0.5rem; display: grid; grid-template-columns: 1fr; gap: 2rem; align-items: center; justify-items: center; height: 100vh;">
		<div style="display: grid; grid-template-columns: 1fr; background-color: white; width: 100%%; border-radius: 0.5rem; align-items: center; gap: 2rem; padding: 2rem; margin: auto; text-align: center;">
			<div>
				<img style="height: auto; width: 10rem; object-fit: contain; margin-top: 2rem;" src="https://everhealth-asset.irfancen.com/assets/eh.png" alt="Everhealth logo" />
			</div>

			<div style="width: 25rem; margin: auto;">
				<h1 style="color: black; margin: 0;">Hello!</h1>
				<p>Password Reset <span style="color: #36A5B2; font-weight: bold;">Everhealth!</span></p>
				<br />
				<p>We received a request to reset your password. Click the button below to reset it:</p>
				<br />
				<a href="%s" style="text-decoration: none; color: white; background-color: #36A5B2; padding: 10px 20px; border-radius: 0.3rem; font-weight: bold; display: inline-block;">Reset Password</a>
				<br />
				<p>If you didn't request a password reset, please ignore this email. Your password will remain unchanged.</p>
				<p>Best,</p>
				<p style="margin-bottom: 3rem; color: #36A5B2; font-weight: bold;">Everhealth</p>
			</div>
		</div>
	</div>
	`, link)
	return subject, content
}

func (r *smtpGmail) SendEmail(token, to string, isVerify bool) error {
	receiver := []string{to}
	link := r.prefixLink + token
	var subject, content string
	if isVerify {
		subject, content = emailVerifyContent(link)
	} else {
		subject, content = emailForgotPasswordContent(link)
	}
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", r.name, r.address)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = receiver
	smtpAuth := smtp.PlainAuth("", r.address, r.password, smtpAuthAddress)
	return e.Send(smtpServerAddress, smtpAuth)
}

func (r *smtpGmail) SendEmailTest(token, to string, isVerify bool) error {
	receiver := []string{to}
	link := r.prefixLink + token
	var subject, content string
	if isVerify {
		subject, content = emailVerifyContent(link)
	} else {
		subject, content = emailForgotPasswordContent(link)
	}
	message := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, content)
	smtpAuth := smtp.PlainAuth("", "", "", "localhost")
	return smtp.SendMail(smtpTestServer, smtpAuth, r.address, receiver, []byte(message))
}

package services

import "log"

// EmailService is a temporary implementation.
// Replace this with Mailtrap / SES later.
type EmailService struct{}

// NewEmailService creates a temp email service
func NewEmailService() *EmailService {
	return &EmailService{}
}

// SendUserInvite simulates sending onboarding email
func (e *EmailService) SendUserInvite(
	to string,
	username string,
	tempPassword string,
) error {
	log.Printf(
		"[DEV EMAIL] To=%s Username=%s TempPassword=%s",
		to,
		username,
		tempPassword,
	)
	return nil
}

// SendPasswordReset simulates password reset email
func (e *EmailService) SendPasswordReset(
	to string,
	resetLink string,
) error {
	log.Printf(
		"[DEV EMAIL] To=%s ResetLink=%s",
		to,
		resetLink,
	)
	return nil
}

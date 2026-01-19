package services

import (
	"github.com/pquerna/otp/totp"
	
)

type OTPService struct{}

func NewOTPService() *OTPService {
	return &OTPService{}
}

func (o *OTPService) GenerateSecret(email string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "RBAC-Auth",
		AccountName: email,
	})
	if err != nil {
		return "", "", err
	}
	return key.Secret(), key.URL(), nil
}

func (o *OTPService) Verify(secret, code string) bool {
	return totp.Validate(code, secret)
}

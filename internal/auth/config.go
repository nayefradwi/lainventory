package auth

import "os"

type AuthConfig struct {
	OtpSecret string
}

func NewAuthConfig() AuthConfig {
	return AuthConfig{
		OtpSecret: os.Getenv("OTP_SECRET"),
	}
}

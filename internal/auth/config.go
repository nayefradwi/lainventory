package auth

import "os"

type AuthConfig struct {
	OtpSecret string
}

func NewAuthConfig() AuthConfig {
	otpSecret := os.Getenv("OTP_SECRET")
	if otpSecret == "" {
		panic("OTP_SECRET is required")
	}

	return AuthConfig{
		OtpSecret: otpSecret,
	}
}

package config

import "testing"

func setRequiredEnv(t *testing.T) {
	t.Setenv("GOOSE_DBSTRING", "postgres://db")
	t.Setenv("REDIS_URL", "redis://cache")
	t.Setenv("OTP_SECRET", "secret")
}

func TestNewConfigSuccess(t *testing.T) {
	setRequiredEnv(t)

	c := NewConfig()

	if c.DatabaseUrl != "postgres://db" || c.RedisUrl != "redis://cache" || c.Auth.OtpSecret != "secret" {
		t.Fatalf("unexpected config: %+v", c)
	}
}

func TestNewConfigPanicsOnMissingEnv(t *testing.T) {
	for _, key := range []string{"GOOSE_DBSTRING", "REDIS_URL", "OTP_SECRET"} {
		t.Run(key, func(t *testing.T) {
			setRequiredEnv(t)
			t.Setenv(key, "")
			defer func() {
				if recover() == nil {
					t.Fatalf("expected panic when %s is empty", key)
				}
			}()
			NewConfig()
		})
	}
}

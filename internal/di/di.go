package di

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	appauth "github.com/nayefradwi/lainventory/internal/auth"
	"github.com/nayefradwi/lainventory/internal/config"
	"github.com/nayefradwi/nayef_go_common/auth"
	"github.com/nayefradwi/nayef_go_common/auth/authhttp"
	"github.com/nayefradwi/nayef_go_common/auth/authpg"
	"github.com/nayefradwi/nayef_go_common/errors"
	"github.com/nayefradwi/nayef_go_common/locking"
	"github.com/nayefradwi/nayef_go_common/pgutil"
	"github.com/nayefradwi/nayef_go_common/redisutil"
	"github.com/redis/go-redis/v9"
)

type Di struct {
	Pool   *pgxpool.Pool
	Redis  *redis.Client
	Locker locking.ILocker
	Auth   AppAuth
}

type AppAuth struct {
	Service      appauth.AuthService
	Authenticate func(http.Handler) http.Handler
	RateLimit    func(http.Handler) http.Handler
}

func RegisterServices(ctx context.Context, config config.Config) *Di {
	di := Di{}
	di.connectToDb(ctx, config.DatabaseUrl)
	di.connectToRedis(ctx, config.RedisUrl)
	di.Locker = locking.NewDistributedLockerFromClient(di.Redis)
	di.registerAuth(config.Auth)
	return &di
}

func (d *Di) Dispose() {
	d.Pool.Close()
	d.Redis.Close()
}

func (d *Di) connectToDb(ctx context.Context, connectionUrl string) {
	d.Pool = pgutil.ConnectToPostgres(ctx, connectionUrl)
}

func (d *Di) connectToRedis(ctx context.Context, connectionUrl string) {
	d.Redis = redisutil.ConnectToRedis(ctx, connectionUrl)
}

func (d *Di) registerAuth(config appauth.AuthConfig) {
	sessionStore := errors.Must(authpg.NewSessionStore(d.Pool))
	codeStore := errors.Must(authpg.NewCodeStore(d.Pool))
	attemptStore := authpg.NewAttemptStore(d.Pool)
	sessions := errors.Must(auth.NewOpaqueSessionManager(sessionStore, 15*time.Minute, 30*24*time.Hour))
	codes := errors.Must(auth.NewOtpManager(codeStore, auth.OtpConfig{
		Secret:     []byte(config.OtpSecret),
		Length:     6,
		CodeLimits: auth.CodeLimits{TTL: 10 * time.Minute, ResendAfter: time.Minute, MaxSends: 5, MaxAttempts: 5},
	}))

	d.Auth = AppAuth{
		Service:      appauth.NewAuthService(d.Pool, sessions, codes),
		Authenticate: authhttp.NewOpaqueMiddleware(sessions).UseAuthentication,
		RateLimit:    authhttp.RateLimit(errors.Must(auth.NewLimiter(attemptStore, 20, time.Hour)), authhttp.ByIP("codes")),
	}
}

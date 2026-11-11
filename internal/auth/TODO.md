# Auth: what's left

Do the steps in order.

## 1. Config

In `internal/config/config.go`, import `github.com/nayefradwi/lainventory/internal/auth` and add:

```go
// add to Config
Auth auth.AuthConfig

// add to NewConfig()
Auth: auth.NewAuthConfig(),
```

## 2. DI

In `internal/di/di.go`, add to `Di`:

```go
Auth          appauth.AuthService
Authenticate  func(http.Handler) http.Handler // use it on your own routes too
AuthRateLimit func(http.Handler) http.Handler
```

Imports:

```go
"net/http"
"time"
"github.com/nayefradwi/nayef_go_common/auth"
"github.com/nayefradwi/nayef_go_common/auth/authpg"
"github.com/nayefradwi/nayef_go_common/auth/authhttp"
"github.com/nayefradwi/nayef_go_common/errors"
appauth "github.com/nayefradwi/lainventory/internal/auth"
```

Add this to `RegisterServices()`, after the database connection is made. You can change any limit or time below.

```go
sessionStore := errors.Must(authpg.NewSessionStore(di.Pool))
codeStore := errors.Must(authpg.NewCodeStore(di.Pool))
attemptStore := authpg.NewAttemptStore(di.Pool)
sessions := errors.Must(auth.NewOpaqueSessionManager(sessionStore, 15*time.Minute, 30*24*time.Hour))
di.Authenticate = authhttp.NewOpaqueMiddleware(sessions).UseAuthentication
di.AuthRateLimit = authhttp.RateLimit(errors.Must(auth.NewLimiter(attemptStore, 20, time.Hour)), authhttp.ByIP("codes"))
codes := errors.Must(auth.NewOtpManager(codeStore, auth.OtpConfig{
	Secret:     []byte(config.Auth.OtpSecret),
	Length:     6,
	CodeLimits: auth.CodeLimits{TTL: 10 * time.Minute, ResendAfter: time.Minute, MaxSends: 5, MaxAttempts: 5},
}))
di.Auth = appauth.NewAuthService(di.Pool, sessions, codes)
```

## 3. Router

In `cmd/api/router.go`, import `appauth "github.com/nayefradwi/lainventory/internal/auth"`. Add to `setup()`:

```go
ar.mux.Mount("/auth", appauth.NewAuthHandler(ar.di.Auth).Routes(ar.di.Authenticate, ar.di.AuthRateLimit))

// protect your own routes
ar.mux.Group(func(r chi.Router) {
	r.Use(ar.di.Authenticate)
	r.Use(authhttp.VerifyClaim("role", "admin")) // optional
})
```

## 4. Secrets

Your local `.env` already has values for `OTP_SECRET`.

Every deployed environment needs its own values. **Do not copy the values from your local `.env`.** Make a new one for each key:

```sh
openssl rand -hex 32
```

Where to put them: open `deployments/<env>/TODO.md` and follow "Add your own environment variables". Treat each key as a secret.

## 5. Database

Run `goose up` on your computer. Deployed environments run it on every deploy.

Expired sessions, codes and attempts stay in the database. To clean them up, run `DeleteExpired` on the session, code and attempt stores on a schedule.

## 6. TODOs

Some parts are left for you. Find them with `grep -rn TODO internal/auth`:

- send login codes
- put roles or scopes in claims

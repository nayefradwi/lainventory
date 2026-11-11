package auth

import (
	"net/http"

	"github.com/nayefradwi/nayef_go_common/errors"
)

var (
	ErrNotLoggedIn     = errors.UnauthorizedError("not logged in")
	ErrUserNotFound    = errors.NotFoundError("user not found")
	ErrTooManyRequests = errors.NewResultErrorWithStatus("Too many requests", "TOO_MANY_REQUESTS", http.StatusTooManyRequests)
)

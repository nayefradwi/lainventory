package auth

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nayefradwi/lainventory/internal/infra/sqlc/gen"
	"github.com/nayefradwi/nayef_go_common/auth"
	"github.com/nayefradwi/nayef_go_common/errors"
	"github.com/nayefradwi/nayef_go_common/pgutil"
)

type AuthService struct {
	q        *db.Queries
	sessions auth.SessionManager
	codes    auth.OtpManager
}

func NewAuthService(
	pool *pgxpool.Pool,
	sessions auth.SessionManager,
	codes auth.OtpManager,
) AuthService {
	return AuthService{
		q:        db.New(pool),
		sessions: sessions,
		codes:    codes,
	}
}

// TODO: put the user's role or scopes in claims
func (svc AuthService) session(ctx context.Context, user db.User) (auth.TokenDTO, error) {
	return svc.sessions.Issue(ctx, pgutil.UUIDFromPg(user.ID), nil)
}

func (svc AuthService) refresh(ctx context.Context, req RefreshRequest) (auth.TokenDTO, error) {
	if err := req.Validate(); err != nil {
		return auth.TokenDTO{}, err
	}

	return svc.sessions.Refresh(ctx, req.RefreshToken)
}

func (svc AuthService) logout(ctx context.Context, req RefreshRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	return svc.sessions.Revoke(ctx, req.RefreshToken)
}

func (svc AuthService) me(ctx context.Context) (UserResponse, error) {
	id, ok := auth.GetIdentity(ctx)
	if !ok {
		return UserResponse{}, ErrNotLoggedIn
	}

	user, err := svc.q.GetUserByID(ctx, pgutil.UUID(id.OwnerId))
	if errors.Is(err, pgx.ErrNoRows) {
		return UserResponse{}, ErrUserNotFound
	}

	if err != nil {
		return UserResponse{}, err
	}

	return UserResponse{Id: id.OwnerId, Email: user.Email}, nil
}

func (svc AuthService) logoutAll(ctx context.Context) error {
	id, ok := auth.GetIdentity(ctx)
	if !ok {
		return ErrNotLoggedIn
	}

	return svc.sessions.RevokeOwner(ctx, id.OwnerId)
}

func (svc AuthService) requestCode(ctx context.Context, req CodeRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	code, err := svc.codes.Issue(ctx, "login:"+req.Email)
	if errors.Is(err, auth.ErrCodeResendBlocked) {
		return ErrTooManyRequests
	}

	if err != nil {
		return err
	}

	// TODO: send the code to req.Email
	_ = code
	return nil
}

func (svc AuthService) verifyCode(ctx context.Context, req VerifyCodeRequest) (auth.TokenDTO, error) {
	if err := req.Validate(); err != nil {
		return auth.TokenDTO{}, err
	}

	if err := svc.codes.Verify(ctx, "login:"+req.Email, req.Code); err != nil {
		return auth.TokenDTO{}, err
	}

	user, err := svc.q.FindOrCreateUser(ctx, req.Email)
	if err != nil {
		return auth.TokenDTO{}, err
	}

	return svc.session(ctx, user)
}

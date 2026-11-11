package auth

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nayefradwi/nayef_go_common/httputil"
)

type AuthHandler struct {
	svc AuthService
}

func NewAuthHandler(svc AuthService) AuthHandler {
	return AuthHandler{
		svc: svc,
	}
}

func (h AuthHandler) Routes(authenticate, rateLimit func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()
	r.Post("/refresh", h.refresh)
	r.Post("/logout", h.logout)
	r.With(rateLimit).Post("/otp/request", h.requestCode)
	r.Post("/otp/verify", h.verifyCode)
	r.Group(func(r chi.Router) {
		r.Use(authenticate)
		r.Get("/me", h.me)
		r.Post("/logout-all", h.logoutAll)
	})

	return r
}

func (h AuthHandler) refresh(w http.ResponseWriter, r *http.Request) {
	httputil.ParseJsonBody(w, r.Body, func(jw httputil.JsonResponseWriter, req RefreshRequest) {
		tokens, err := h.svc.refresh(r.Context(), req)
		jw.WriteJsonResponse(tokens, err)
	})
}

func (h AuthHandler) logout(w http.ResponseWriter, r *http.Request) {
	httputil.ParseJsonBody(w, r.Body, func(jw httputil.JsonResponseWriter, req RefreshRequest) {
		err := h.svc.logout(r.Context(), req)
		jw.WriteSuccessMessage("logged out", err)
	})
}

func (h AuthHandler) me(w http.ResponseWriter, r *http.Request) {
	jw := httputil.NewJsonResponseWriter(w)
	user, err := h.svc.me(r.Context())
	jw.WriteJsonResponse(user, err)
}

func (h AuthHandler) logoutAll(w http.ResponseWriter, r *http.Request) {
	jw := httputil.NewJsonResponseWriter(w)
	err := h.svc.logoutAll(r.Context())
	jw.WriteSuccessMessage("logged out everywhere", err)
}

func (h AuthHandler) requestCode(w http.ResponseWriter, r *http.Request) {
	httputil.ParseJsonBody(w, r.Body, func(jw httputil.JsonResponseWriter, req CodeRequest) {
		err := h.svc.requestCode(r.Context(), req)
		jw.WithSuccessStatus(http.StatusAccepted).WriteSuccessMessage("code sent", err)
	})
}

func (h AuthHandler) verifyCode(w http.ResponseWriter, r *http.Request) {
	httputil.ParseJsonBody(w, r.Body, func(jw httputil.JsonResponseWriter, req VerifyCodeRequest) {
		res, err := h.svc.verifyCode(r.Context(), req)
		jw.WriteJsonResponse(res, err)
	})
}

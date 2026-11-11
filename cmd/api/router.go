package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nayefradwi/lainventory/internal/di"
	"github.com/nayefradwi/lainventory/internal/health"
	"github.com/nayefradwi/nayef_go_common/pagination"
)

type router struct {
	mux    *chi.Mux
	server *http.Server
	di     *di.Di
}

func newRouter(address string, di *di.Di) router {
	mux := chi.NewRouter()
	r := router{
		mux: mux,
		di:  di,
		server: &http.Server{
			Addr:    address,
			Handler: mux,
		},
	}

	r.setup()
	return r
}

func (ar router) setup() {
	ar.mux.Use(middleware.Logger)
	ar.mux.Use(middleware.RequestID)
	ar.mux.Use(middleware.StripSlashes)
	ar.mux.Use(pagination.SetOffsetPaginationMiddleware)
	ar.mux.Get("/health", health.HealthCheck)
}

func (ar router) connect(cancel context.CancelFunc) {
	slog.Info("staring server on address", "address", ar.server.Addr)
	if err := ar.server.ListenAndServe(); err != nil {
		slog.Error("failed to start server", "error", err.Error())
		cancel()
	}
}

func (ar router) disconnect(ctx context.Context) {
	slog.Info("disconnecting server", "address", ar.server.Addr)
	if err := ar.server.Shutdown(ctx); err != nil {
		slog.Error("failed to stop server", "error", err.Error())
		panic(err)
	}
}

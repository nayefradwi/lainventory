package main

import (
	"context"
	"github.com/nayefradwi/lainventory/internal/config"
	"github.com/nayefradwi/lainventory/internal/di"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func bootstrap(ctx context.Context) {
	conf := config.NewConfig()

	d := di.RegisterServices(ctx, conf)
	defer d.Dispose()

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	router := newRouter(conf.Address, d)
	go router.connect(stop)
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	router.disconnect(shutdownCtx)
}

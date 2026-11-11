package di

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nayefradwi/lainventory/internal/config"
	"github.com/nayefradwi/nayef_go_common/locking"
	"github.com/nayefradwi/nayef_go_common/pgutil"
	"github.com/nayefradwi/nayef_go_common/redisutil"
	"github.com/redis/go-redis/v9"
)

type Di struct {
	Pool   *pgxpool.Pool
	Redis  *redis.Client
	Locker locking.ILocker
}

func RegisterServices(ctx context.Context, config config.Config) *Di {
	di := Di{}
	di.connectToDb(ctx, config.DatabaseUrl)
	di.connectToRedis(ctx, config.RedisUrl)
	di.Locker = locking.NewDistributedLockerFromClient(di.Redis)
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

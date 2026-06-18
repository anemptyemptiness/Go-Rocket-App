package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	apiauth "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/api/auth/v1"
	apiuser "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/api/user/v1"
	"github.com/anemptyemptiness/Go-Rocket-App/iam/internal/config"
	sessionrepo "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/repository/session"
	userrepo "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/repository/user"
	authsvc "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/service/auth"
	usersvc "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/service/user"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/closer"
	pkgredis "github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/redis"
	authv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/auth/v1"
	userv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/user/v1"
)

type diContainer struct {
	redisClient *redis.Client
	pgPool      *pgxpool.Pool

	authApi authv1.AuthServiceServer
	userApi userv1.UserServiceServer

	authService apiauth.AuthService
	userService apiuser.UserService

	sessionRepository authsvc.SessionRepository
	userRepository    usersvc.UserRepository
}

func (d *diContainer) PgPool(ctx context.Context) *pgxpool.Pool {
	if d.pgPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().PG.DSN())
		if err != nil {
			slog.Error("создание пула соединений", "error", err)
			os.Exit(1)
		}

		closer.Add("pgPool", func(_ context.Context) error {
			pool.Close()
			return nil
		})

		err = pool.Ping(ctx)
		if err != nil {
			slog.Error("проверка соединения с БД", "error", err)
			os.Exit(1)
		}

		slog.Info("подключение к PostgreSQL установлено")

		d.pgPool = pool
	}

	return d.pgPool
}

func (d *diContainer) RedisClient(ctx context.Context) *redis.Client {
	if d.redisClient == nil {
		client, err := pkgredis.NewClient(
			&redis.Options{
				Addr:            config.AppConfig().Redis.Address(),
				DialTimeout:     config.AppConfig().Redis.ConnectionTimeout,
				ReadTimeout:     config.AppConfig().Redis.ConnectionTimeout,
				WriteTimeout:    config.AppConfig().Redis.ConnectionTimeout,
				MaxIdleConns:    config.AppConfig().Redis.MaxIdle,
				ConnMaxIdleTime: config.AppConfig().Redis.IdleTimeout,
			},
			slog.Default(),
		)
		if err != nil {
			slog.Error("ошибка подключения к Redis", "error", err)
			os.Exit(1)
		}

		closer.Add("redis client", func(_ context.Context) error {
			return client.Close()
		})

		d.redisClient = client
	}

	return d.redisClient
}

func (d *diContainer) AuthAPI(ctx context.Context) authv1.AuthServiceServer {
	if d.authApi == nil {
		d.authApi = apiauth.New(d.AuthService(ctx))
	}

	return d.authApi
}

func (d *diContainer) UserAPI(ctx context.Context) userv1.UserServiceServer {
	if d.userApi == nil {
		d.userApi = apiuser.New(d.UserService(ctx))
	}

	return d.userApi
}

func (d *diContainer) AuthService(ctx context.Context) apiauth.AuthService {
	if d.authService == nil {
		d.authService = authsvc.New(d.SessionRepository(ctx), d.UserRepository(ctx))
	}

	return d.authService
}

func (d *diContainer) UserService(ctx context.Context) apiuser.UserService {
	if d.userService == nil {
		d.userService = usersvc.New(d.UserRepository(ctx))
	}

	return d.userService
}

func (d *diContainer) SessionRepository(ctx context.Context) authsvc.SessionRepository {
	if d.sessionRepository == nil {
		d.sessionRepository = sessionrepo.New(d.RedisClient(ctx))
	}

	return d.sessionRepository
}

func (d *diContainer) UserRepository(ctx context.Context) usersvc.UserRepository {
	if d.userRepository == nil {
		d.userRepository = userrepo.New(d.PgPool(ctx))
	}

	return d.userRepository
}

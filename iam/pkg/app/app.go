package app

import (
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	apiauth "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/api/auth/v1"
	apiuser "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/api/user/v1"
	sessionrepo "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/repository/session"
	userrepo "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/repository/user"
	authsvc "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/service/auth"
	usersvc "github.com/anemptyemptiness/Go-Rocket-App/iam/internal/service/user"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
	authv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/auth/v1"
	userv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/user/v1"
)

const (
	maxConnectionIdle     = 15 * time.Minute
	maxConnectionAge      = 30 * time.Minute
	maxConnectionAgeGrace = 5 * time.Second

	keepAliveTime                = 5 * time.Minute
	keepAliveTimeout             = 10 * time.Second
	keepAliveMinTime             = 50 * time.Second
	keepAlivePermitWithoutStream = true
)

func NewGRPCServer(
	pool *pgxpool.Pool,
	redisDB *redis.Client,
	ttl time.Duration,
	cryptCost int,
) *grpc.Server {
	userRepo := userrepo.New(pool)
	sessionRepo := sessionrepo.New(redisDB, ttl)

	userSvc := usersvc.New(userRepo, cryptCost)
	authSvc := authsvc.New(sessionRepo, userRepo)

	userApi := apiuser.New(userSvc)
	authApi := apiauth.New(authSvc)

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     maxConnectionIdle,
			MaxConnectionAge:      maxConnectionAge,
			MaxConnectionAgeGrace: maxConnectionAgeGrace,
			Time:                  keepAliveTime,
			Timeout:               keepAliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             keepAliveMinTime,
			PermitWithoutStream: keepAlivePermitWithoutStream,
		}),
		grpc.UnaryInterceptor(pkgerr.UnaryErrorInterceptor(slog.Default())),
	)

	userv1.RegisterUserServiceServer(grpcServer, userApi)
	authv1.RegisterAuthServiceServer(grpcServer, authApi)

	return grpcServer
}

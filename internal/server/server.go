package server

import (
	"github.com/AleksK1NG/auth-microservice/config"
	"github.com/AleksK1NG/auth-microservice/internal/interceptors"
	authServerGRPC "github.com/AleksK1NG/auth-microservice/internal/user/delivery/grpc/server"
	"github.com/AleksK1NG/auth-microservice/internal/user/repository"
	"github.com/AleksK1NG/auth-microservice/internal/user/usecase"
	"github.com/AleksK1NG/auth-microservice/pkg/logger"
	userService "github.com/AleksK1NG/auth-microservice/proto"
	"github.com/go-redis/redis/v8"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

// GRPC Auth Server
type Server struct {
	logger      logger.Logger
	cfg         *config.Config
	db          *sqlx.DB
	redisClient *redis.Client
}

// Server constructor
func NewAuthServer(logger logger.Logger, cfg *config.Config, db *sqlx.DB, redisClient *redis.Client) *Server {
	return &Server{logger: logger, cfg: cfg, db: db, redisClient: redisClient}
}

// Run server
func (s *Server) Run() error {
	im := interceptors.NewInterceptorManager(s.logger, s.cfg)
	userRepo := repository.NewUserRepository(s.db)
	userUC := usecase.NewUserUseCase(s.logger, userRepo)

	l, err := net.Listen("tcp", s.cfg.Server.Port)
	if err != nil {
		return err
	}

	server := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: 5 * time.Minute,
		Timeout:           15 * time.Second,
		MaxConnectionAge:  5 * time.Minute,
	}),
		grpc.UnaryInterceptor(im.Logger),
		grpc.ChainUnaryInterceptor(grpcrecovery.UnaryServerInterceptor()),
	)

	if s.cfg.Server.Mode != "Production" {
		reflection.Register(server)
	}

	authGRPCServer := authServerGRPC.NewAuthServerGRPC(s.logger, s.cfg, userUC)
	userService.RegisterUserServiceServer(server, authGRPCServer)

	s.logger.Infof("Server is listening on port: %v", s.cfg.Server.Port)
	if err := server.Serve(l); err != nil {
		return err
	}
	return nil
}

package server

import (
	"github.com/AleksK1NG/auth-microservice/config"
	"github.com/AleksK1NG/auth-microservice/internal/interceptors"
	sessRepository "github.com/AleksK1NG/auth-microservice/internal/session/repository"
	sessUseCase "github.com/AleksK1NG/auth-microservice/internal/session/usecase"
	authServerGRPC "github.com/AleksK1NG/auth-microservice/internal/user/delivery/grpc/service"
	userRepository "github.com/AleksK1NG/auth-microservice/internal/user/repository"
	userUseCase "github.com/AleksK1NG/auth-microservice/internal/user/usecase"
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

// Run service
func (s *Server) Run() error {
	im := interceptors.NewInterceptorManager(s.logger, s.cfg)
	userRepo := userRepository.NewUserRepository(s.db)
	sessRepo := sessRepository.NewSessionRepository(s.redisClient, s.cfg)
	userRedisRepo := userRepository.NewUserRedisRepo(s.redisClient, s.logger)
	userUC := userUseCase.NewUserUseCase(s.logger, userRepo, userRedisRepo)
	sessUC := sessUseCase.NewSessionUseCase(sessRepo, s.cfg)

	l, err := net.Listen("tcp", s.cfg.Server.Port)
	if err != nil {
		return err
	}

	server := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: s.cfg.Server.MaxConnectionIdle * time.Minute,
		Timeout:           s.cfg.Server.Timeout * time.Second,
		MaxConnectionAge:  s.cfg.Server.MaxConnectionAge * time.Minute,
		Time:              s.cfg.Server.Timeout * time.Minute,
	}),
		grpc.UnaryInterceptor(im.Logger),
		grpc.ChainUnaryInterceptor(grpcrecovery.UnaryServerInterceptor()),
	)

	if s.cfg.Server.Mode != "Production" {
		reflection.Register(server)
	}

	authGRPCServer := authServerGRPC.NewAuthServerGRPC(s.logger, s.cfg, userUC, sessUC)
	userService.RegisterUserServiceServer(server, authGRPCServer)

	s.logger.Infof("Server is listening on port: %v", s.cfg.Server.Port)
	if err := server.Serve(l); err != nil {
		return err
	}
	return nil
}

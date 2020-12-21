package server

import (
	"github.com/AleksK1NG/auth-microservice/config"
	"github.com/AleksK1NG/auth-microservice/internal/interceptors"
	authServerGRPC "github.com/AleksK1NG/auth-microservice/internal/user/delivery/grpc/server"
	"github.com/AleksK1NG/auth-microservice/pkg/logger"
	userService "github.com/AleksK1NG/auth-microservice/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

// GRPC Auth Server
type Server struct {
	logger logger.Logger
	cfg    *config.Config
}

// Server constructor
func NewAuthServer(logger logger.Logger, cfg *config.Config) *Server {
	return &Server{logger: logger, cfg: cfg}
}

func (s *Server) Run() error {
	im := interceptors.NewInterceptorManager(s.logger, s.cfg)

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
	)

	if s.cfg.Server.Mode != "Production" {
		reflection.Register(server)
	}

	authGRPCServer := authServerGRPC.NewAuthServerGRPC(s.logger, s.cfg)
	userService.RegisterUserServiceServer(server, authGRPCServer)

	s.logger.Infof("Server is listening on port: %v", s.cfg.Server.Port)
	if err := server.Serve(l); err != nil {
		return err
	}
	return nil
}

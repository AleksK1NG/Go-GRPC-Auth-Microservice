package interceptors

import (
	"context"
	"github.com/AleksK1NG/auth-microservice/config"
	"github.com/AleksK1NG/auth-microservice/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

// InterceptorManager
type InterceptorManager struct {
	logger logger.Logger
	cfg    *config.Config
}

// InterceptorManager constructor
func NewInterceptorManager(logger logger.Logger, cfg *config.Config) *InterceptorManager {
	return &InterceptorManager{logger: logger, cfg: cfg}
}

// Logger Interceptor
func (im *InterceptorManager) Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	reply, err := handler(ctx, req)
	im.logger.Infof("METHOD: %s, reply: %#v, time: %v, metadata: %v, err: %#v", info.FullMethod, reply, time.Since(start), md, err)

	return reply, err
}

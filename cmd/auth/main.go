package main

import (
	"log"
	"os"

	"github.com/opentracing/opentracing-go"

	"github.com/AleksK1NG/auth-microservice/config"
	"github.com/AleksK1NG/auth-microservice/internal/server"
	jaegerTracer "github.com/AleksK1NG/auth-microservice/pkg/jaeger"
	"github.com/AleksK1NG/auth-microservice/pkg/logger"
	"github.com/AleksK1NG/auth-microservice/pkg/postgres"
	"github.com/AleksK1NG/auth-microservice/pkg/redis"
	"github.com/AleksK1NG/auth-microservice/pkg/utils"
)

func main() {
	log.Println("Starting auth microservice")

	configPath := utils.GetConfigPath(os.Getenv("config"))
	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	appLogger := logger.NewAPILogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof(
		"AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v",
		cfg.Server.AppVersion,
		cfg.Logger.Level,
		cfg.Server.Mode,
		cfg.Server.SSL,
	)
	appLogger.Infof("Success parsed config: %#v", cfg.Server.AppVersion)

	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	}
	defer psqlDB.Close()

	redisClient := redis.NewRedisClient(cfg)
	defer redisClient.Close()
	appLogger.Info("Redis connected")

	tracer, closer, err := jaegerTracer.InitJaeger(cfg)
	if err != nil {
		appLogger.Fatal("cannot create tracer", err)
	}
	appLogger.Info("Jaeger connected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	appLogger.Info("Opentracing connected")

	authServer := server.NewAuthServer(appLogger, cfg, psqlDB, redisClient)
	appLogger.Fatal(authServer.Run())
}

package main

import (
	"log"
	"net"
	"os"
	"time"

	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/rezaAmiri123/service-user/app/handler"
	"github.com/rezaAmiri123/service-user/app/interceptors"
	"github.com/rezaAmiri123/service-user/app/model"
	"github.com/rezaAmiri123/service-user/app/repository"
	"github.com/rezaAmiri123/service-user/cmd/config"
	pb "github.com/rezaAmiri123/service-user/gen/pb"
	"github.com/rezaAmiri123/service-user/pkg/jaeger"
	"github.com/rezaAmiri123/service-user/pkg/logger"
	"github.com/rezaAmiri123/service-user/pkg/metric"
	"github.com/rezaAmiri123/service-user/pkg/mysql"
	"github.com/rezaAmiri123/service-user/pkg/utils"
)

func main() {
	log.Println("Starting user server microservice")

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

	db := mysql.NewGormDB(cfg)
	defer db.Close()
	model.AutoMigrate(db)

	metrics, err := metric.CreateMetrics(cfg.Metrics.URL, cfg.Metrics.ServiceName)
	if err != nil {
		appLogger.Errorf("CreateMetrics Error: %s", err)
	}
	appLogger.Info(
		"Metrics available URL: %s, ServiceName: %s",
		cfg.Metrics.URL,
		cfg.Metrics.ServiceName,
	)

	im := interceptors.NewInterceptorManager(appLogger, cfg, metrics)

	tracer, closer, err := jaeger.InitJaeger(cfg)
	if err != nil {
		appLogger.Fatal("cannot create tracer", err)
	}
	appLogger.Info("Jaeger connected")
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	appLogger.Info("Opentracing connected")

	repo := repository.NewORMUserRepository(db)
	h := handler.NewUserHandler(repo, appLogger)
	lis, err := net.Listen("tcp", cfg.Server.Port)
	if err != nil {
		appLogger.Fatal(err.Error())
	}

	srv := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: cfg.Server.MaxConnectionIdle * time.Minute,
		Timeout:           cfg.Server.Timeout * time.Second,
		MaxConnectionAge:  cfg.Server.MaxConnectionAge * time.Minute,
		Time:              cfg.Server.Timeout * time.Minute,
	}),
		grpc.UnaryInterceptor(im.Logger),
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpcrecovery.UnaryServerInterceptor(),
		),
	)

	pb.RegisterUsersServer(srv, h)
	appLogger.Info("server starts at ", cfg.Server.Port)

	if err := srv.Serve(lis); err != nil {
		appLogger.Fatal(err.Error())
	}

}

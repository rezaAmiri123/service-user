package main

import (
	"context"
	"net/http"

	"github.com/alexflint/go-arg"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rezaAmiri123/service-user/cmd/config"
	pb "github.com/rezaAmiri123/service-user/gen/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func run(cfg *config.Config) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ropts := []runtime.ServeMuxOption{
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{}),
	}

	mux := runtime.NewServeMux(ropts...)
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := pb.RegisterUsersHandlerFromEndpoint(context.Background(), mux, cfg.GetServerAddress(), opts)
	if err != nil {
		return err
	}
	logrus.Printf("starting gateway server on port %v", cfg.GatewayPort)
	return http.ListenAndServe(cfg.GetGatewayAddress(), mux)
}

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	cfg := config.DefaultConfiguration()
	arg.MustParse(cfg)

	if err := run(cfg); err != nil {
		logrus.Fatal(err.Error())
	}
}

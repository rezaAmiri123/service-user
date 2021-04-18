package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/rezaAmiri123/service-user/cmd/config"
	pb "github.com/rezaAmiri123/service-user/gen/pb"
	"github.com/rezaAmiri123/service-user/pkg/trace"
	"github.com/rezaAmiri123/service-user/pkg/utils"
)

func run(cfg *config.Config) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ropts := []runtime.ServeMuxOption{
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{}),
	}

	mux := runtime.NewServeMux(ropts...)
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			grpc_opentracing.UnaryClientInterceptor(
				grpc_opentracing.WithTracer(opentracing.GlobalTracer()),
			),
		),
	}

	err := pb.RegisterUsersHandlerFromEndpoint(context.Background(), mux, cfg.Gateway.GetServerAddress(), opts)
	if err != nil {
		return err
	}
	mux.HandlePath("GET", "/swagger.json", serveSwagger)

	mux.HandlePath("GET", "/swagger-ui", serveSwaggerFiles)
	logrus.Printf("starting gateway server on port %v", cfg.Gateway.Port)
	return http.ListenAndServe(cfg.Gateway.Port, trace.TracingWrapper(mux))
}

func main() {
	log.Println("Starting user gateway microservice")

	configPath := utils.GetConfigPath(os.Getenv("config"))
	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	if err := run(cfg); err != nil {
		logrus.Fatal(err.Error())
	}
}

func serveSwagger(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	http.ServeFile(w, r, "gen/swagger/user.swagger.json")
}

func serveSwaggerFiles(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	// TODO this function dosen't work
	fs := http.FileServer(http.Dir("templates/swagger-ui"))
	http.StripPrefix("/static/", fs)
	//http.ServeFile()
	//http.StripPrefix()
	//http.Handle()
	//fs.ServeHTTP(w, r)
}

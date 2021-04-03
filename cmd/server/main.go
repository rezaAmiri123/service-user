package main

import (
	"net"

	"github.com/alexflint/go-arg"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/rezaAmiri123/service-user/app/handler"
	"github.com/rezaAmiri123/service-user/app/model"
	"github.com/rezaAmiri123/service-user/app/repository"
	"github.com/rezaAmiri123/service-user/cmd/config"
	pb "github.com/rezaAmiri123/service-user/gen/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	// Initialize config struct and populate it froms env vars and flags.
	cfg := config.DefaultConfiguration()
	arg.MustParse(cfg)

	db := config.SetupDB(cfg)
	model.AutoMigrate(db)
	repo := repository.NewORMUserRepository(db)
	h := handler.NewUserHandler(repo)
	lis, err := net.Listen("tcp", cfg.GetAddress())
	if err != nil {
		logrus.Fatal(err.Error())
	}

	srv := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(),
		),
	)

	pb.RegisterUsersServer(srv, h)
	logrus.Infoln("server starts at ", cfg.GetAddress())

	if err := srv.Serve(lis); err != nil {
		logrus.Fatal(err.Error())
	}

}

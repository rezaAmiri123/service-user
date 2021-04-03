package handler

import (
	"context"
	pb "github.com/rezaAmiri123/service-user/gen/pb"
)

type UserHandler struct{}

func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {

}

func (h *UserHandler) GetUser(ctx context.Context, req *pb.Empty) (*pb.UserResponse, error) {

}

func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {

}

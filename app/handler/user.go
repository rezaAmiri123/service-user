package handler

import (
	"context"
	"fmt"
	"github.com/rezaAmiri123/service-user/app/auth"
	"github.com/rezaAmiri123/service-user/app/model"
	"github.com/rezaAmiri123/service-user/app/repository"
	pb "github.com/rezaAmiri123/service-user/gen/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler{
	return &UserHandler{repo: repo}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	u := &model.User{
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
	if err := u.Validate(); err != nil {
		msg := fmt.Sprintf("validation error: %w", err.Error())
		return nil, status.Error(codes.InvalidArgument, msg)
	}
	if err := u.HashPassword(); err != nil {
		msg := fmt.Sprintf("failed to hash password: %w", err.Error())
		return nil, status.Error(codes.Aborted, msg)
	}
	if err := h.repo.Create(u); err != nil {
		msg := fmt.Sprintf("failed to create user: %w", err.Error())
		return nil, status.Error(codes.Canceled, msg)
	}
	return u.ProtoResponse(), nil
}

// LoginUser is existing user login
func (h *UserHandler) LoginUser(ctx context.Context,req *pb.LoginRequest) (*pb.LoginResponse, error){
	user, err := h.repo.GetByEmail(req.GetEmail())
	if err != nil{
		msg := fmt.Sprintf("invalid email or password: %w", err.Error())
		return nil, status.Error(codes.InvalidArgument, msg)
	}
	if !user.CheckPassword(req.GetPassword()){
		msg := fmt.Sprintf("invalid email or password: %w", err.Error())
		return nil, status.Error(codes.InvalidArgument, msg)
	}
	token,err := auth.GenerateToken(user.ID)
	if err != nil{
		msg := fmt.Sprintf("failed to create token: %w", err.Error())
		return nil, status.Error(codes.InvalidArgument, msg)
	}
	return &pb.LoginResponse{Token: token},nil
}

// GetUser gets current user
func (h *UserHandler) GetUser(ctx context.Context, req *pb.Empty) (*pb.UserResponse, error) {
	u, err := h.getUser(ctx)
	if err != nil{
		return nil, err
	}
	return u.ProtoResponse(), nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	u, err := h.getUser(ctx)
	if err != nil{
		return nil, err
	}
	// update non zero-valu fields eonly
	username := req.GetUsername()
	if username != "" {
		u.Username = username
	}

	email := req.GetEmail()
	if email != "" {
		u.Email = email
	}

	password := req.GetPassword()
	if password != "" {
		u.Password = password
		u.HashPassword()
	}

	if err := u.Validate();err != nil{
		msg := fmt.Sprintf("validation: %w", err.Error())
		return nil, status.Error(codes.InvalidArgument, msg)
	}

	if err := h.repo.Update(u);err!= nil{
		msg := fmt.Sprintf("failed to update: %w", err.Error())
		return nil, status.Error(codes.InvalidArgument, msg)
	}
	return u.ProtoResponse(), nil
}

func (h *UserHandler) getUser(ctx context.Context) (*model.User, error) {
	userID,err := auth.GetUserID(ctx)
	if err != nil{
		msg := fmt.Sprintf("unauthenticated: %w", err.Error())
		return nil, status.Error(codes.Unauthenticated, msg)
	}
	u ,err := h.repo.GetByID(userID)
	if err != nil{
		msg := fmt.Sprintf("user not found: %w", err.Error())
		return nil, status.Error(codes.NotFound, msg)
	}

	return u, nil
}

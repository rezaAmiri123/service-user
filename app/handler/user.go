package handler

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/rezaAmiri123/service-user/app/auth"
	"github.com/rezaAmiri123/service-user/app/model"
	"github.com/rezaAmiri123/service-user/app/repository"
	pb "github.com/rezaAmiri123/service-user/gen/pb"
	"github.com/rezaAmiri123/service-user/pkg/logger"
)

type UserHandler struct {
	repo   repository.UserRepository
	logger logger.Logger
}

func NewUserHandler(repo repository.UserRepository, logger logger.Logger) *UserHandler {
	return &UserHandler{repo: repo, logger: logger}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Create")
	defer span.Finish()

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
	if err := h.repo.Create(ctx, u); err != nil {
		msg := fmt.Sprintf("failed to create user: %w", err.Error())
		return nil, status.Error(codes.Canceled, msg)
	}
	return u.ProtoUser(), nil
}

// LoginUser is existing user login
func (h *UserHandler) LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Login")
	defer span.Finish()

	user, err := h.repo.GetByEmail(ctx, req.GetEmail())
	if err != nil {
		msg := fmt.Sprintf("invalid email or password: %w", err.Error())
		return nil, status.Error(codes.InvalidArgument, msg)
	}
	if !user.CheckPassword(req.GetPassword()) {
		msg := fmt.Sprintf("invalid email or password: %w", err.Error())
		return nil, status.Error(codes.InvalidArgument, msg)
	}
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		msg := fmt.Sprintf("failed to create token: %w", err.Error())
		return nil, status.Error(codes.InvalidArgument, msg)
	}
	return &pb.LoginResponse{Token: token}, nil
}

// GetUser gets current user
func (h *UserHandler) GetUser(ctx context.Context, req *pb.Empty) (*pb.UserResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.GetUser")
	defer span.Finish()

	u, err := h.getUser(ctx)
	if err != nil {
		return nil, err
	}
	return u.ProtoUser(), nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.UpdateUser")
	defer span.Finish()

	u, err := h.getUser(ctx)
	if err != nil {
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

	if err := u.Validate(); err != nil {
		msg := fmt.Sprintf("validation: %w", err.Error())
		return nil, status.Error(codes.InvalidArgument, msg)
	}

	if err := h.repo.Update(ctx, u); err != nil {
		msg := fmt.Sprintf("failed to update: %w", err.Error())
		return nil, status.Error(codes.InvalidArgument, msg)
	}
	return u.ProtoUser(), nil
}

func (h *UserHandler) GetProfile(ctx context.Context, req *pb.ProfileRequest) (*pb.ProfileResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.GetProfile")
	defer span.Finish()

	u, err := h.getUser(ctx)
	if err != nil {
		return nil, err
	}
	otherUser, err := h.repo.GetByUsername(ctx, req.GetUsername())
	if err != nil {
		msg := fmt.Sprintf("user not found: %w", err.Error())
		return nil, status.Error(codes.NotFound, msg)
	}
	isFollowing, err := h.repo.IsFollowing(ctx, u, otherUser)
	if err != nil {
		msg := fmt.Sprintf("failed to get follow status: %w", err.Error())
		return nil, status.Error(codes.NotFound, msg)
	}
	return u.ProtoProfile(isFollowing), nil
}

func (h *UserHandler) FollowUser(ctx context.Context, req *pb.FollowRequest) (*pb.ProfileResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.FollowUser")
	defer span.Finish()

	u, err := h.getUser(ctx)
	if err != nil {
		return nil, err
	}
	if u.Username == req.GetUsername() {
		msg := fmt.Sprintf("cannot follow yourself: %w", err.Error())
		return nil, status.Error(codes.InvalidArgument, msg)
	}
	otherUser, err := h.repo.GetByUsername(ctx, req.GetUsername())
	if err != nil {
		msg := fmt.Sprintf("user not found: %w", err.Error())
		return nil, status.Error(codes.NotFound, msg)
	}
	if err := h.repo.Follow(ctx, u, otherUser); err != nil {
		msg := fmt.Sprintf("failed to follow user: %w", err.Error())
		return nil, status.Error(codes.NotFound, msg)
	}
	return u.ProtoProfile(true), nil
}

func (h *UserHandler) UnFollowUser(ctx context.Context, req *pb.FollowRequest) (*pb.ProfileResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.UnfollowUser")
	defer span.Finish()

	u, err := h.getUser(ctx)
	if err != nil {
		return nil, err
	}
	if u.Username == req.GetUsername() {
		msg := fmt.Sprintf("cannot unfollow yourself: %w", err.Error())
		return nil, status.Error(codes.InvalidArgument, msg)
	}
	otherUser, err := h.repo.GetByUsername(ctx, req.GetUsername())
	if err != nil {
		msg := fmt.Sprintf("user not found: %w", err.Error())
		return nil, status.Error(codes.NotFound, msg)
	}
	following, err := h.repo.IsFollowing(ctx, u, otherUser)
	if err != nil {
		msg := fmt.Sprintf("failed to get following: %w", err.Error())
		return nil, status.Error(codes.NotFound, msg)
	}
	if !following {
		msg := fmt.Sprintf("user is not following the other user: %w", err.Error())
		return nil, status.Error(codes.NotFound, msg)
	}
	if err := h.repo.Unfollow(ctx, u, otherUser); err != nil {
		msg := fmt.Sprintf("failed to unfollow user: %w", err.Error())
		return nil, status.Error(codes.Aborted, msg)
	}
	return u.ProtoProfile(false), nil
}

func (h *UserHandler) getUser(ctx context.Context) (*model.User, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		msg := fmt.Sprintf("unauthenticated: %w", err.Error())
		return nil, status.Error(codes.Unauthenticated, msg)
	}
	u, err := h.repo.GetByID(ctx, userID)
	if err != nil {
		msg := fmt.Sprintf("user not found: %w", err.Error())
		return nil, status.Error(codes.NotFound, msg)
	}
	return u, nil
}

package grpc

import (
	"context"

	"github.com/KritSuth/fullstack-coding-test/backend-go/internal/model"
	"github.com/KritSuth/fullstack-coding-test/backend-go/internal/service"
	pb "github.com/KritSuth/fullstack-coding-test/backend-go/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserGRPCServer implements the gRPC UserService interface.
type UserGRPCServer struct {
	pb.UnimplementedUserServiceServer
	svc *service.UserService
}

func NewUserGRPCServer(svc *service.UserService) *UserGRPCServer {
	return &UserGRPCServer{svc: svc}
}

func (s *UserGRPCServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	createReq := &model.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	if err := createReq.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	user, err := s.svc.Register(ctx, createReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.UserResponse{
		Id:        user.ID.Hex(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
	}, nil
}

func (s *UserGRPCServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	user, err := s.svc.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &pb.UserResponse{
		Id:        user.ID.Hex(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
	}, nil
}

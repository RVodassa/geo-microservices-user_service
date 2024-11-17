package userGrpcService

import (
	"context"
	"geo-microservices/user/internal/domain/entity"

	"geo-microservices/user/internal/service"
	proto "geo-microservices/user/proto/generated"
)

type UserServiceServer struct {
	proto.UnimplementedUserServiceServer
	service service.UserServiceProvider
}

func NewUserServiceServer(service service.UserServiceProvider) *UserServiceServer {
	return &UserServiceServer{service: service}
}

func (s *UserServiceServer) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	newUser := &entity.User{
		Password: req.Password,
		Login:    req.Login,
	}
	id, err := s.service.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return &proto.CreateUserResponse{Id: id}, nil
}

func (s *UserServiceServer) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	err := s.service.DeleteUser(ctx, uint64(req.Id))
	if err != nil {
		return nil, err
	}

	return &proto.DeleteUserResponse{Success: true}, nil
}

package service

import (
	"context"
	"geo-microservices/user/internal/domain/entity"
	"geo-microservices/user/internal/repository"
)

type UserServiceProvider interface {
	CreateUser(ctx context.Context, u *entity.User) (id uint64, err error)
	DeleteUser(ctx context.Context, id uint64) error
}

type UserService struct {
	repo repository.UserRepositoryProvider
}

func NewUserService(repo repository.UserRepositoryProvider) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, u *entity.User) (id uint64, err error) {
	return s.repo.Register(ctx, u)
}

func (s *UserService) DeleteUser(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}

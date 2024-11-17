package service

import (
	"context"
	"github.com/RVodassa/geo-microservices-user/internal/domain/entity"
	"github.com/RVodassa/geo-microservices-user/internal/repository"
)

type UserServiceProvider interface {
	Register(ctx context.Context, u *entity.User) (id uint64, err error)
	Login(ctx context.Context, login, password string) (bool, error)
	Delete(ctx context.Context, id uint64) error
	Profile(ctx context.Context, id uint64) (*entity.User, error)
	List(ctx context.Context, offset, limit uint64) ([]*entity.User, uint64, error)
}

type UserService struct {
	repo repository.UserRepositoryProvider
}

func NewUserService(repo repository.UserRepositoryProvider) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, u *entity.User) (id uint64, err error) {
	return s.repo.Register(ctx, u)
}

func (s *UserService) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}

func (s *UserService) Profile(ctx context.Context, id uint64) (*entity.User, error) {
	return s.repo.Profile(ctx, id)
}
func (s *UserService) List(ctx context.Context, offset, limit uint64) ([]*entity.User, uint64, error) {
	return s.repo.List(ctx, offset, limit)
}

func (s *UserService) Login(ctx context.Context, login, password string) (bool, error) {
	return s.repo.Login(ctx, login, password)
}

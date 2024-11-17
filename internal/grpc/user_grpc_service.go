package userGrpcServer

import (
	"context"
	"github.com/RVodassa/geo-microservices-user/internal/domain/entity"
	"github.com/RVodassa/geo-microservices-user/internal/service"
	pb "github.com/RVodassa/geo-microservices-user/proto/generated"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"sync"
)

type User struct {
	ID        uint64
	Login     string
	Password  string // Храним хэш пароля
	CreatedAt string
}

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	mu      sync.Mutex
	service service.UserServiceProvider
}

func NewUserServiceServer(service service.UserServiceProvider) *UserServiceServer {
	return &UserServiceServer{
		service: service,
	}
}

// Register добавляет нового пользователя.
func (s *UserServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := entity.User{
		Login:    req.Login,
		Password: req.Password, // В реальной системе здесь должен быть bcrypt-хэш.
	}
	id, err := s.service.Register(ctx, &user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Printf("Registered user: %+v", user)
	return &pb.RegisterResponse{Id: id}, nil
}

// Delete удаляет пользователя по ID.
func (s *UserServiceServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := entity.User{
		ID: req.Id,
	}

	err := s.service.Delete(ctx, user.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Printf("Deleted user with ID: %d", req.Id)
	return &pb.DeleteResponse{Status: true}, nil
}

// Profile возвращает профиль пользователя по ID.
func (s *UserServiceServer) Profile(ctx context.Context, req *pb.ProfileRequest) (*pb.ProfileResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, err := s.service.Profile(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ProfileResponse{
		Id:    user.ID,
		Login: user.Login,
	}, nil
}

// List возвращает список пользователей с пагинацией.
func (s *UserServiceServer) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Получение данных через сервисный слой
	users, count, err := s.service.List(ctx, req.Offset, req.Limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch user list: %v", err)
	}

	// Инициализация ответа
	userList := make([]*pb.User, 0, len(users)) // Предварительное выделение памяти
	for _, user := range users {
		userList = append(userList, &pb.User{
			Id:    user.ID,
			Login: user.Login,
		})
	}

	return &pb.ListResponse{
		Users: userList,
		Total: count,
	}, nil
}

func (s *UserServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	exist, err := s.service.Login(ctx, req.Login, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.LoginResponse{
		Status: exist,
	}, nil
}

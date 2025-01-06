package main

import (
	grpcservice "github.com/RVodassa/geo-microservices-user_service/internal/grpc-server"
	"github.com/RVodassa/geo-microservices-user_service/internal/repository"
	"github.com/RVodassa/geo-microservices-user_service/internal/service"
	"github.com/RVodassa/geo-microservices-user_service/internal/sql"
	proto "github.com/RVodassa/geo-microservices-user_service/proto/generated"
	"google.golang.org/grpc"
	"net"
	"sync"
	"time"
)

func main() {

	time.Sleep(3 * time.Second)

	db, err := sql.InitDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Создание сервера
	grpcServer := grpc.NewServer()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	userGrpcService := grpcservice.NewUserServiceServer(userService)

	proto.RegisterUserServiceServer(grpcServer, userGrpcService)

	var wg sync.WaitGroup
	errChan := make(chan error)

	wg.Add(1)
	go func() {
		defer wg.Done()

		listener, err := net.Listen("tcp", ":10101")
		if err != nil {
			errChan <- err

			return
		}

		if err := grpcServer.Serve(listener); err != nil {
			errChan <- err

			return
		}
	}()

	if err := <-errChan; err != nil {
		panic(err)
	}

	wg.Wait()
}

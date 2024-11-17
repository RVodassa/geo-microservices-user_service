package main

import (
	grpc_service "geo-microservices/user/internal/grpc"
	"geo-microservices/user/internal/repository"
	"geo-microservices/user/internal/service"
	"geo-microservices/user/internal/sql"
	proto "geo-microservices/user/proto/generated"
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

	userGrpcService := grpc_service.NewUserServiceServer(userService)

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

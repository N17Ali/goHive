package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/n17ali/gohive/api/taskpb"
	"github.com/n17ali/gohive/internal/logs"
	"github.com/n17ali/gohive/internal/storage/redis"
	"github.com/n17ali/gohive/internal/task"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	redis.InitRedis("localhost:6379", 0)

	ln, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	logger := logs.RedisLogger{}
	grpcServer := grpc.NewServer()
	taskService := task.NewTaskService(logger)
	taskpb.RegisterTaskServiceServer(grpcServer, taskService)

	reflection.Register(grpcServer)

	ctx := context.Background()
	executor := task.NewTaskExecutor(10*time.Second, logger)
	go executor.Start(ctx)

	fmt.Println("🚀 gRPC Server running on port 50051...")
	if err := grpcServer.Serve(ln); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

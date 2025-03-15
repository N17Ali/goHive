package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/n17ali/gohive/internal/pb"
	"github.com/n17ali/gohive/internal/task"
	"github.com/n17ali/gohive/pkg/redis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	redis.InitRedis("localhost:6379", 0)

	ln, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTaskServiceServer(grpcServer, &task.TaskServiceServer{})

	reflection.Register(grpcServer)

	ctx := context.Background()
	executor := task.NewTaskExecutor(10 * time.Second)
	go executor.Start(ctx)

	fmt.Println("🚀 gRPC Server running on port 50051...")
	if err := grpcServer.Serve(ln); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

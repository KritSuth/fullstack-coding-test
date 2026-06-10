package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	grpcserver "github.com/KritSuth/fullstack-coding-test/backend-go/internal/grpc"
	"github.com/KritSuth/fullstack-coding-test/backend-go/internal/repository"
	"github.com/KritSuth/fullstack-coding-test/backend-go/internal/service"
	pb "github.com/KritSuth/fullstack-coding-test/backend-go/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("mongo connect: %v", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("userapi")
	repo := repository.NewMongoUserRepository(db)
	svc := service.NewUserService(repo)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, grpcserver.NewUserGRPCServer(svc))

	log.Println("gRPC server running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

package main

import (
	"log"
	"net"

	"gRPC-Learner-Project/proto/logger"

	"google.golang.org/grpc"
)

type Server struct {
	logger.UnimplementedLoggerServer
}

func main() {

	if err := InitLogger(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer Logger.Sync()

	listen, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal("Failed to listen on Port 9000: ", err)
	}

	grpcServer := grpc.NewServer()

	logger.RegisterLoggerServer(grpcServer, &Server{})

	log.Println("Logger Service is starting on port 9000...")

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatal("Failed to start logger service on port 9000: ", err)
	}
}

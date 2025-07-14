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
	listen, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal("Failed to listen on Port 9000: ", err)
	}

	grpcServer := grpc.NewServer()

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatal("Failed to start logger service on port 9000: ", err)
	}

	logger.RegisterLoggerServer(grpcServer, &Server{})

	log.Println("Logger Service is running on port 9000")
}

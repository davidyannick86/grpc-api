package main

import (
	"log"
	"net"
	"os"

	"github.com/davidyannick86/grpc-api-mongodb/internals/api/handlers"
	"github.com/davidyannick86/grpc-api-mongodb/internals/api/interceptors"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/davidyannick86/grpc-api-mongodb/proto/gen"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors.ResponseTimeInterceptor))

	pb.RegisterExecsServiceServer(server, &handlers.Server{})
	pb.RegisterStudentsServiceServer(server, &handlers.Server{})
	pb.RegisterTeachersServiceServer(server, &handlers.Server{})

	port := os.Getenv("SERVER_PORT")

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Server is running on port", port)

	reflection.Register(server)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

package main

import (
	"context"
	"log"
	"net"

	"github.com/Tobi696/blog-application/proto"
	"google.golang.org/grpc"
)

type authServer struct{}

func (authServer) Login(_ context.Context, in *proto.LoginRequest) (*proto.AuthResponse, error) {
	return &proto.AuthResponse{}, nil
}
func main() {

	server := grpc.NewServer()
	proto.RegisterAuthServiceServer(server, authServer)

	listener, err := net.Listen("tcp", "5000")
	if err != nil {
		log.Fatal("Error creating listener : ", err.Error())
	}

	server.Serve(listener)
}

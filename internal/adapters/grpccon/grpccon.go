package grpccon

import (
	"errors"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userPb "github.com/LordMoMA/Hexagonal-Architecture/internal/core/proto/user"
)

func NewGrpcClient(host string) (userPb.UserGrpcServiceClient, error) {
	opts := make([]grpc.DialOption, 0)

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	clientConn, err := grpc.NewClient(host, opts...)
	if err != nil {
		log.Printf("Error: Grpc client connection failed: %s", err.Error())
		return nil, errors.New("error: grpc client connection failed")
	}

	return userPb.NewUserGrpcServiceClient(clientConn), nil
}

func NewGrpcServer(host string) (*grpc.Server, net.Listener) {
	opts := make([]grpc.ServerOption, 0)

	grpcServer := grpc.NewServer(opts...)

	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	return grpcServer, lis
}

package grpc

import (
	"fmt"
	pb "github.com/Royal17x/flagr/backend/gen/proto/v1"
	"github.com/Royal17x/flagr/backend/internal/port"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func NewGRPCServer(flagSvc port.FlagServiceInterface, authSvc port.AuthServiceInterface) (*grpc.Server, *FlagServer) {
	srv := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			UnaryLoggingInterceptor,
			UnaryAuthInterceptor(authSvc),
		),
	)

	flagServer := NewFlagServer(flagSvc)
	pb.RegisterFlagServiceServer(srv, flagServer)

	reflection.Register(srv)
	return srv, flagServer
}

func StartGRPCServer(srv *grpc.Server, port string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return fmt.Errorf("grpc: listen: %w", err)
	}
	return srv.Serve(lis)
}

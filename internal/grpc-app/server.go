package grpc_app

import (
	"context"
	"net"

	"google.golang.org/grpc"
)

type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

type server struct {
	grpcServer *grpc.Server
	address    string
}

func (s *server) ListenAndServe() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	return s.grpcServer.Serve(listener)
}

func (s *server) Shutdown(ctx context.Context) error {
	s.grpcServer.GracefulStop()
	return nil
}

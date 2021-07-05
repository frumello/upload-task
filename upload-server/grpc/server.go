package grpc

import (
	"log"
	"net"
	"os"
	"time"

	"upload-server/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Server struct {
	lis    net.Listener
	server *grpc.Server
}

func New(grpcTarget, httpTarget string) *Server {
	lis, err := net.Listen("tcp", grpcTarget)
	if err != nil {
		os.Exit(2)
	}

	handler := NewHandler(httpTarget)
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionAge: 2 * time.Minute}),
	)
	pb.RegisterImageServiceServer(grpcServer, handler)

	return &Server{
		lis:    lis,
		server: grpcServer,
	}
}

func (s *Server) ListenAndServe() error {
	log.Printf("GRPC server serving at : %s", s.lis.Addr())
	return s.server.Serve(s.lis)
}

func (s *Server) GracefulStop() {
	log.Printf("GRPC server shutting down")
	s.server.GracefulStop()
}

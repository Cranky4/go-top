package internalgrpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/Cranky4/go-top/api/TopService"
	"github.com/Cranky4/go-top/internal/top"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedTopServiceServer
	grpcServer        *grpc.Server
	logg              Logger
	app               *top.Top
	requestLogFile    string
	requestLogHandler *os.File
}

func New(app *top.Top, logg Logger, requestLogFile string) *Server {
	return &Server{app: app, logg: logg, requestLogFile: requestLogFile}
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Error(msg string)
}

func (s *Server) Start(ctx context.Context, addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	opts := []grpc.ServerOption{}
	s.grpcServer = grpc.NewServer(opts...)

	file, err := os.Create(s.requestLogFile)
	if err != nil {
		return err
	}
	logger := log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)

	handler, err := NewHandler(s.app, logger)
	if err != nil {
		return err
	}

	pb.RegisterTopServiceServer(s.grpcServer, handler)
	go s.grpcServer.Serve(listener)
	s.logg.Info(fmt.Sprintf("grpc server started and listen %s...", addr))

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) {
	s.grpcServer.Stop()
	s.requestLogHandler.Close()
	s.logg.Info("grpc server stopped")
}

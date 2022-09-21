package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Cranky4/go-top/internal/logger"
	internalgrpc "github.com/Cranky4/go-top/internal/server/grpc"
	"github.com/Cranky4/go-top/internal/top"
)

var grpcAddr, configFile string

func init() {
	flag.StringVar(&grpcAddr, "grpc-addr", ":9990", "GRPC server port")
	flag.StringVar(&configFile, "config", "./config/app.toml", "Path to config file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	config := NewConfig(configFile)
	config.Grpc.Addr = grpcAddr
	logg := logger.New(config.Logg.Level, log.LstdFlags)

	top := top.NewTop(config, *logg)

	var grpcServer *internalgrpc.Server
	go func() {
		grpcServer = internalgrpc.New(top, logg, config.Grpc.RequestLogFile)
		if err := grpcServer.Start(ctx, config.Grpc.Addr); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			cancel()
			os.Exit(1)
		}

		logg.Info("grpc started on " + config.Grpc.Addr)
	}()
	defer cancel()

	logg.Info("top is running...")

	<-ctx.Done()
}

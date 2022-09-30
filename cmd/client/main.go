package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Cranky4/go-top/internal/logger"
	topclient "github.com/Cranky4/go-top/internal/top-client"
)

var (
	grpcAddr, configFile string
	m, n                 int
)

func init() {
	flag.StringVar(&grpcAddr, "grpc-addr", ":9990", "GRPC server port")
	flag.IntVar(&m, "m", 15, "Snapshot period seconds")
	flag.IntVar(&n, "n", 5, "Snapshot offset seconds")
	flag.StringVar(&configFile, "config", "./config/client.toml", "Path to config file")
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
	config.Client.M = m
	config.Client.N = n

	logg := logger.New(config.Logg.Level, log.LstdFlags)
	client := topclient.New(ctx, config, *logg)

	if err := client.Start(); err != nil {
		logg.Error("failed to start top client: " + err.Error())
		cancel()
		os.Exit(1)
	}
	defer cancel()

	logg.Info("client is running...")

	<-ctx.Done()
}

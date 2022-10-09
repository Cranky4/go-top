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
	grpcAddr, configFile          string
	warmingUpTime, snapshotPeriod int
)

func init() {
	flag.StringVar(&grpcAddr, "grpc-addr", ":9990", "GRPC server port")
	flag.IntVar(&warmingUpTime, "m", 15, "Snapshot warming up time (seconds)")
	flag.IntVar(&snapshotPeriod, "n", 5, "Snapshot period (seconds)")
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
	config.Client.WarmingUpTime = warmingUpTime
	config.Client.SnapshotPeriod = snapshotPeriod

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

	logg.Info("client is stopped...")
}

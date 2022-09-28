package app

import (
	appdf "github.com/Cranky4/go-top/internal/app/df"
	appiostat "github.com/Cranky4/go-top/internal/app/iostat"
	appnetstat "github.com/Cranky4/go-top/internal/app/netstat"
	apptcpdump "github.com/Cranky4/go-top/internal/app/tcpdump"
	apptop "github.com/Cranky4/go-top/internal/app/top"
)

type Config struct {
	App     AppConf
	Metrics MetricsConf
	Logg    LoggerConf
	Grpc    GrpcConf
}

type AppConf struct {
	TopPath, IostatPath, DfPath, TcpDumpPath, TimeoutPath, NetStatPath string
}
type MetricsConf struct {
	Cpu, Disks, Network, Connections bool
}

type LoggerConf struct {
	Level string
}

type GrpcConf struct {
	Addr, RequestLogFile string
}

type Snapshot struct {
	Cpu                  apptop.Cpu
	DisksIO              []appiostat.DiskIO
	DisksInfo            []appdf.DiskInfo
	TopTalkersByProtocol []apptcpdump.TopTalkerByProtocol
	TopTalkersByTraffic  []apptcpdump.TopTalkerByTraffic
	ConnectsInfo         []appnetstat.ConnectInfo
	ConnectsStates       []appnetstat.ConnectState
}

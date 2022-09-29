package app

import (
	appdf "github.com/Cranky4/go-top/internal/app/df"
	appiostat "github.com/Cranky4/go-top/internal/app/iostat"
	appnetstat "github.com/Cranky4/go-top/internal/app/netstat"
	apptcpdump "github.com/Cranky4/go-top/internal/app/tcpdump"
	apptop "github.com/Cranky4/go-top/internal/app/top"
)

type Config struct {
	App     Conf
	Metrics MetricsConf
	Logg    LoggerConf
	Grpc    GrpcConf
}

type Conf struct {
	TopPath, IostatPath, DfPath, TCPDumpPath, TimeoutPath, NetStatPath string
}
type MetricsConf struct {
	CPU, Disks, Network, Connections bool
}

type LoggerConf struct {
	Level string
}

type GrpcConf struct {
	Addr, RequestLogFile string
}

type Snapshot struct {
	CPU                  apptop.CPU
	DisksIO              []appiostat.DiskIO
	DisksInfo            []appdf.DiskInfo
	TopTalkersByProtocol []apptcpdump.TopTalkerByProtocol
	TopTalkersByTraffic  []apptcpdump.TopTalkerByTraffic
	ConnectsInfo         []appnetstat.ConnectInfo
	ConnectsStates       []appnetstat.ConnectState
}

package app

import (
	"time"

	appdf "github.com/Cranky4/go-top/internal/app/df"
	appiostat "github.com/Cranky4/go-top/internal/app/iostat"
	apptop "github.com/Cranky4/go-top/internal/app/top"
)

type Config struct {
	App  AppConf
	Logg LoggerConf
	Grpc GrpcConf
}

type AppConf struct {
	Metrics, TopPath, IostatPath, DfPath string
}

type LoggerConf struct {
	Level string
}

type GrpcConf struct {
	Addr, RequestLogFile string
}

type Snapshot struct {
	StartTime, FinishTime time.Time
	Cpu                   apptop.Cpu
	DisksIO               []appiostat.DiskIO
	DisksInfo             []appdf.DiskInfo
	TopTalkersByProtocol  []TopTalkerByProtocol
	TopTalkersByTraffic   []TopTalkerByTraffic
	ConnectsInfo          []ConnectInfo
	ConnectsStates        []ConnectState
}

type ConnectInfo struct {
	Command  string // -
	Pid      int32  // -
	User     string // ?
	Protocol string // TCP
	Port     string // 40349
}

type ConnectState struct {
	Protocol string // tcp
	State    string // listen
}

type TopTalkerByProtocol struct {
	Protocol string // UDP
	Bytes    int64  // 127
	Percent  string // 32%
}

type TopTalkerByTraffic struct {
	Source         string // 172.21.0.1.52978
	Bytes          int64  // 239.255.255.250.1900
	Protocol       string // udp
	BytesPerSecond int64  // 173 ?
}

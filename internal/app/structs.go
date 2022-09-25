package app

import (
	"time"

	apptop "github.com/Cranky4/go-top/internal/app/top"
)

type Config struct {
	Top  TopConf
	Logg LoggerConf
	Grpc GrpcConf
}

type TopConf struct {
	Metrics, TopPath string
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
	DisksIO               []DiskIO
	DisksInfo             []DiskInfo
	TopTalkersByProtocol  []TopTalkerByProtocol
	TopTalkersByTraffic   []TopTalkerByTraffic
	ConnectsInfo          []ConnectInfo
	ConnectsStates        []ConnectState
}

type DiskIO struct {
	Device, Tps, Kbps string // nvme0n1,  52.86, 665.63 + 780.25
}

type DiskInfo struct {
	Name            string // /dev/nvme0n1p3     df -k
	UsedBytes       int64  // 39131452
	AvailableBytes  int64  // 64169664
	UsageBytes      string // 38%
	UsedInodes      int64  // 272911         df -i
	AvailableInodes int64  // 64224433
	UsageInodes     string // 1%
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

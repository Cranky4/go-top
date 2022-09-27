package apptcpdump

import "time"

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type TopTalkers struct {
	ByProtocol TopTalkerByProtocol
	ByTraffic  TopTalkerByTraffic
}

type TcpDumpLine struct {
	Time                          time.Time
	Protocol, Source, Destination string
	Bytes                         int64
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

type Parser interface {
	Parse(in string) ([]TcpDumpLine, error)
}

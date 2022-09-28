package apptcpdump

import (
	"fmt"
	"time"
)

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type TopTalkers struct {
	ByProtocol []TopTalkerByProtocol
	ByTraffic  []TopTalkerByTraffic
}

type TcpDumpLine struct {
	Time                                time.Time
	Type, Protocol, Source, Destination string
	Bytes                               int
}

type TopTalkerByProtocol struct {
	Protocol string // UDP
	Bytes    int    // 127
	Percent  string // 32%
}

type TopTalkerByTraffic struct {
	Source         string  // 172.21.0.1.52978
	Destination    string  // 239.255.255.250.1900
	Protocol       string  // udp
	BytesPerSecond float32 // 173 ?
}

type Parser interface {
	Parse(in string) ([]TcpDumpLine, error)
}

// errors
type ErrCannotParseInput struct {
	Input string
}

func (e *ErrCannotParseInput) Error() string {
	return fmt.Sprintf("cannot parse %s", e.Input)
}

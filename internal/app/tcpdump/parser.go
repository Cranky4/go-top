package apptcpdump

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

type TcpDumpParser struct {
	dataReg *regexp.Regexp
}

func NewParser() *TcpDumpParser {
	return &TcpDumpParser{
		dataReg: regexp.MustCompile(
			`(\d+:\d+:\d+.\d+).*?(\d+.\d+.\d+.\d+.\d+).*?>.*?(\d+.\d+.\d+.\d+.\d+).*?(tcp|udp|icmp).*?[length]?\s(\d+)$`,
		),
	}
}

func (t *TcpDumpParser) Parse(in string) ([]TcpDumpLine, error) {
	rows := strings.Split(in, "\n")
	result := make([]TcpDumpLine, 0, len(rows))

	for i := 1; i < len(rows); i++ {
		if rows[i] == "" {
			continue
		}

		parts := t.dataReg.FindStringSubmatch(rows[i])

		time, err := time.Parse("15:04:05.999999999", parts[1])
		if err != nil {
			return nil, err
		}

		source := parts[2]
		destination := parts[3]
		protocol := parts[4]

		bytes, err := strconv.Atoi(parts[5])
		if err != nil {
			return nil, err
		}

		result = append(result, TcpDumpLine{
			Time:        time,
			Source:      source,
			Destination: destination,
			Protocol:    protocol,
			Bytes:       int64(bytes),
		})
	}

	return result, nil
}

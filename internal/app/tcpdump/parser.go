package apptcpdump

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

type TcpDumpParser struct {
	dataReg *regexp.Regexp
	logg    Logger
}

func NewParser(logg Logger) *TcpDumpParser {
	return &TcpDumpParser{
		logg: logg,
		dataReg: regexp.MustCompile(
			`^(\d+\-\d+\-\d+\s\d+:\d+:\d+.\d+).*?(\w+).*?([\w\d\.\:]+).*?>.*?([\w\d\.\:]+)\:.*?(\w+).*?[length]?\s(\d+)$`,
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

		if len(parts) == 0 {
			err := &ErrCannotParseInput{Input: rows[i]}
			t.logg.Warn(err.Error())
			continue
		}

		time, err := time.Parse("2006-01-02 15:04:05.999999999", parts[1])
		if err != nil {
			return nil, err
		}

		typ := parts[2]
		source := parts[3]
		destination := parts[4]
		protocol := parts[5]

		bytes, err := strconv.Atoi(parts[6])
		if err != nil {
			return nil, err
		}

		result = append(result, TcpDumpLine{
			Time:        time,
			Type:        typ,
			Source:      source,
			Destination: destination,
			Protocol:    protocol,
			Bytes:       bytes,
		})
	}

	return result, nil
}

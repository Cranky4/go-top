package appiostat

import (
	"regexp"
	"strconv"
	"strings"
)

type IostatParser struct {
	dataReg *regexp.Regexp
}

func NewParser() *IostatParser {
	return &IostatParser{
		dataReg: regexp.MustCompile(`([\w]+)\s+([\d\.]+)\s+([\d\.]+)\s+([\d\.]+)`),
	}
}

func (t *IostatParser) Parse(in string) ([]DiskIO, error) {
	rows := strings.Split(in, "\n")
	discs := make([]DiskIO, 0, len(rows)-3)

	for i := 3; i < len(rows); i++ {
		if rows[i] == "" {
			continue
		}

		parts := t.dataReg.FindStringSubmatch(rows[i])

		device := parts[1]

		tps, err := strconv.ParseFloat(parts[2], 32)
		if err != nil {
			return nil, err
		}

		read, err := strconv.ParseFloat(parts[3], 32)
		if err != nil {
			return nil, err
		}

		write, err := strconv.ParseFloat(parts[4], 32)
		if err != nil {
			return nil, err
		}

		discs = append(discs, DiskIO{
			Device: device,
			Tps:    float32(tps),
			Kbps:   float32(read + write),
		})
	}

	return discs, nil
}

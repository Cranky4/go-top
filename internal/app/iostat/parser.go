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
		dataReg: regexp.MustCompile(
			`^([\w\d]+)\s+([\d\.]+)\s+([\d\.]+)\s+([\d\.]+)\s+([\d\.]+)\s+([\d\.]+)$`,
		),
	}
}

func (t *IostatParser) Parse(in string) ([]IostatRow, error) {
	rows := strings.Split(in, "\n")
	discs := make([]IostatRow, 0, len(rows)-3)

	for i := 3; i < len(rows); i++ {
		if rows[i] == "" {
			continue
		}

		parts := t.dataReg.FindStringSubmatch(rows[i])

		if len(parts) == 0 {
			return nil, &ErrCannotParseInput{Input: rows[i]}
		}

		device := parts[1]

		tps, err := strconv.ParseFloat(parts[2], 32)
		if err != nil {
			return nil, err
		}

		kbpsRead, err := strconv.ParseFloat(parts[3], 32)
		if err != nil {
			return nil, err
		}

		kbpsWrite, err := strconv.ParseFloat(parts[4], 32)
		if err != nil {
			return nil, err
		}

		kbRead, err := strconv.Atoi(parts[5])
		if err != nil {
			return nil, err
		}

		kbWrite, err := strconv.Atoi(parts[6])
		if err != nil {
			return nil, err
		}

		discs = append(discs, IostatRow{
			Device:    device,
			Tps:       float32(tps),
			KbpsRead:  float32(kbpsRead),
			KbpsWrite: float32(kbpsWrite),
			KbRead:    kbRead,
			KbWrite:   kbWrite,
		})
	}

	return discs, nil
}

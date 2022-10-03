package appiostat

import (
	"regexp"
	"strconv"
	"strings"
)

type IostatParser struct {
	dataReg *regexp.Regexp
	logg    Logger
}

func NewParser(logg Logger) *IostatParser {
	return &IostatParser{
		dataReg: regexp.MustCompile(
			`^(.*?)\s+([\d\.]+)\s+([\d\.]+)\s+([\d\.]+)\s+([\d\.]+)\s+([\d\.]+)$`,
		),
		logg: logg,
	}
}

func (t *IostatParser) Parse(in string) []IostatRow {
	rows := strings.Split(in, "\n")
	discs := make([]IostatRow, 0, len(rows)-3)

	for i := 3; i < len(rows); i++ {
		if rows[i] == "" {
			continue
		}

		parts := t.dataReg.FindStringSubmatch(rows[i])

		if len(parts) == 0 {
			err := &ErrCannotParseInput{Input: rows[i]}
			t.logg.Error(err.Error())
			continue
		}

		device := parts[1]

		tps, err := strconv.ParseFloat(parts[2], 32)
		if err != nil {
			t.logg.Error(err.Error())
			continue
		}

		kbpsRead, err := strconv.ParseFloat(parts[3], 32)
		if err != nil {
			t.logg.Error(err.Error())
			continue
		}

		kbpsWrite, err := strconv.ParseFloat(parts[4], 32)
		if err != nil {
			t.logg.Error(err.Error())
			continue
		}

		kbRead, err := strconv.Atoi(parts[5])
		if err != nil {
			t.logg.Error(err.Error())
			continue
		}

		kbWrite, err := strconv.Atoi(parts[6])
		if err != nil {
			t.logg.Error(err.Error())
			continue
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

	return discs
}

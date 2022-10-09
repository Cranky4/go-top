package appdf

import (
	"regexp"
	"strconv"
	"strings"
)

type DfParser struct {
	dataReg *regexp.Regexp
	logg    Logger
}

func NewParser(logg Logger) *DfParser {
	return &DfParser{
		dataReg: regexp.MustCompile(
			`^(.*?)\s+([\d]+)\s+([\d]+)\s+([\d]+)\s+([\d\-]+)%?\s+(.*?)$`,
		),
		logg: logg,
	}
}

func (t *DfParser) ParseBytes(in string) []DiskInfo {
	rows := strings.Split(in, "\n")
	discs := make([]DiskInfo, 0, len(rows)-1)

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

		device := parts[1]

		used, err := strconv.Atoi(parts[3])
		if err != nil {
			t.logg.Error(err.Error())
			continue
		}

		available, err := strconv.Atoi(parts[4])
		if err != nil {
			t.logg.Error(err.Error())
			continue
		}

		var usage int
		if parts[5] != "-" {
			usage, err = strconv.Atoi(parts[5])
			if err != nil {
				t.logg.Error(err.Error())
				continue
			}
		}

		discs = append(discs, DiskInfo{
			Name:           device,
			UsedBytes:      used,
			AvailableBytes: available,
			UsageBytes:     usage,
		})
	}

	return discs
}

func (t *DfParser) ParseInodes(in string) []DiskInfo {
	rows := strings.Split(in, "\n")
	discs := make([]DiskInfo, 0, len(rows)-1)

	for i := 1; i < len(rows); i++ {
		if rows[i] == "" {
			continue
		}

		parts := t.dataReg.FindStringSubmatch(rows[i])
		if len(parts) < 2 {
			continue
		}

		device := parts[1]

		available, err := strconv.Atoi(parts[2])
		if err != nil {
			t.logg.Error(err.Error())
			continue
		}

		used, err := strconv.Atoi(parts[3])
		if err != nil {
			t.logg.Error(err.Error())
			continue
		}

		var usage int
		if parts[5] != "-" {
			usage, err = strconv.Atoi(parts[5])
			if err != nil {
				t.logg.Error(err.Error())
				continue
			}
		}

		discs = append(discs, DiskInfo{
			Name:            device,
			UsedInodes:      used,
			AvailableInodes: available,
			UsageInodes:     usage,
		})
	}

	return discs
}

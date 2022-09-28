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
			`^([\w\/]+)\s+([\d]+)\s+([\d]+)\s+([\d]+)\s+([\d\-]+)%?\s+([\w\/]+)$`,
		),
		logg: logg,
	}
}

func (t *DfParser) ParseBytes(in string) ([]DiskInfo, error) {
	rows := strings.Split(in, "\n")
	discs := make([]DiskInfo, 0, len(rows)-1)

	for i := 1; i < len(rows); i++ {
		if rows[i] == "" {
			continue
		}

		parts := t.dataReg.FindStringSubmatch(rows[i])

		if len(parts) == 0 {
			return nil, &ErrCannotParseInput{Input: rows[i]}
		}

		device := parts[1]

		used, err := strconv.Atoi(parts[3])
		if err != nil {
			return nil, err
		}

		available, err := strconv.Atoi(parts[4])
		if err != nil {
			return nil, err
		}

		var usage int
		if parts[5] != "-" {
			usage, err = strconv.Atoi(parts[5])
			if err != nil {
				return nil, err
			}
		}

		discs = append(discs, DiskInfo{
			Name:           device,
			UsedBytes:      used,
			AvailableBytes: available,
			UsageBytes:     usage,
		})
	}

	return discs, nil
}

func (t *DfParser) ParseInodes(in string) ([]DiskInfo, error) {
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
			return nil, err
		}

		used, err := strconv.Atoi(parts[3])
		if err != nil {
			return nil, err
		}

		var usage int
		if parts[5] != "-" {
			usage, err = strconv.Atoi(parts[5])
			if err != nil {
				return nil, err
			}
		}

		discs = append(discs, DiskInfo{
			Name:            device,
			UsedInodes:      used,
			AvailableInodes: available,
			UsageInodes:     usage,
		})
	}

	return discs, nil
}

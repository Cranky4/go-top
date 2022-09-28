package apptop

import (
	"regexp"
	"strconv"
	"strings"
)

type TopParser struct {
	avgLoadReg, statesReg *regexp.Regexp
}

func NewParser() *TopParser {
	return &TopParser{
		avgLoadReg: regexp.MustCompile(`load average: ([\d.]+), ([\d.]+), ([\d.]+)`),
		statesReg:  regexp.MustCompile(`([\d\.]+)\s+us,.*?([\d\.]+)\s+sy,.*?([\d\.]+)\s+id`),
	}
}

func (t *TopParser) Parse(in string) (Cpu, error) {
	rows := strings.SplitN(in, "\n", 4)
	// avg
	avgParts := t.avgLoadReg.FindStringSubmatch(rows[0])

	if len(avgParts) == 0 {
		return Cpu{}, &ErrCannotParseInput{Input: rows[0]}
	}

	min, err := strconv.ParseFloat(avgParts[1], 32)
	if err != nil {
		return Cpu{}, err
	}

	five, err := strconv.ParseFloat(avgParts[2], 32)
	if err != nil {
		return Cpu{}, err
	}

	fifteen, err := strconv.ParseFloat(avgParts[3], 32)
	if err != nil {
		return Cpu{}, err
	}

	// states
	statesParts := t.statesReg.FindStringSubmatch(rows[2])

	if len(statesParts) == 0 {
		return Cpu{}, &ErrCannotParseInput{Input: rows[2]}
	}

	us, err := strconv.ParseFloat(statesParts[1], 32)
	if err != nil {
		return Cpu{}, err
	}

	sy, err := strconv.ParseFloat(statesParts[2], 32)
	if err != nil {
		return Cpu{}, err
	}

	id, err := strconv.ParseFloat(statesParts[3], 32)
	if err != nil {
		return Cpu{}, err
	}

	return Cpu{
		Avg: CpuAvg{
			Min:     float32(min),
			Five:    float32(five),
			Fifteen: float32(fifteen),
		},
		State: CpuState{
			User:   float32(us),
			System: float32(sy),
			Idle:   float32(id),
		},
	}, nil
}

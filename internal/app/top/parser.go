package apptop

import (
	"regexp"
	"strconv"
	"strings"
)

type TopParser struct {
	avgLoadReg, statesReg *regexp.Regexp
	logg                  Logger
}

func NewParser(logg Logger) *TopParser {
	return &TopParser{
		avgLoadReg: regexp.MustCompile(`load average: ([\d.]+), ([\d.]+), ([\d.]+)`),
		statesReg:  regexp.MustCompile(`([\d\.]+)\s+us,.*?([\d\.]+)\s+sy,.*?([\d\.]+)\s+id`),
		logg:       logg,
	}
}

func (t *TopParser) Parse(in string) CPU {
	rows := strings.SplitN(in, "\n", 4)
	// avg
	avgParts := t.avgLoadReg.FindStringSubmatch(rows[0])

	if len(avgParts) == 0 {
		err := &ErrCannotParseInput{Input: rows[0]}
		t.logg.Warn(err.Error())
		return CPU{}
	}

	min, err := strconv.ParseFloat(avgParts[1], 32)
	if err != nil {
		t.logg.Error(err.Error())
		return CPU{}
	}

	five, err := strconv.ParseFloat(avgParts[2], 32)
	if err != nil {
		t.logg.Error(err.Error())
		return CPU{}
	}

	fifteen, err := strconv.ParseFloat(avgParts[3], 32)
	if err != nil {
		t.logg.Error(err.Error())
		return CPU{}
	}

	// states
	statesParts := t.statesReg.FindStringSubmatch(rows[2])

	if len(statesParts) == 0 {
		err := &ErrCannotParseInput{Input: rows[2]}
		t.logg.Error(err.Error())
		return CPU{}
	}

	us, err := strconv.ParseFloat(statesParts[1], 32)
	if err != nil {
		t.logg.Error(err.Error())
		return CPU{}
	}

	sy, err := strconv.ParseFloat(statesParts[2], 32)
	if err != nil {
		t.logg.Error(err.Error())
		return CPU{}
	}

	id, err := strconv.ParseFloat(statesParts[3], 32)
	if err != nil {
		t.logg.Error(err.Error())
		return CPU{}
	}

	return CPU{
		Avg: CPUAvg{
			Min:     float32(min),
			Five:    float32(five),
			Fifteen: float32(fifteen),
		},
		State: CPUState{
			User:   float32(us),
			System: float32(sy),
			Idle:   float32(id),
		},
	}
}

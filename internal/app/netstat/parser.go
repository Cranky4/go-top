package appnetstat

import (
	"regexp"
	"strconv"
	"strings"
)

type NetStatParser struct {
	dataReg *regexp.Regexp
	logg    Logger
}

func NewParser(logg Logger) Parser {
	return &NetStatParser{
		dataReg: regexp.MustCompile(
			`^([\w]+)\s+([\d]+)\s+([\d]+)\s+([\d\:\.\*]+)\s+([\d\:\.\*]+)\s+([\w]+)\s+([\d\w\/\-]+)\s+$`,
		),
		logg: logg,
	}
}

func (t *NetStatParser) Parse(in string) []NetStatRow {
	rows := strings.Split(in, "\n")
	connects := make([]NetStatRow, 0, len(rows)-1)

	for i := 2; i < len(rows); i++ {
		if rows[i] == "" {
			continue
		}

		parts := t.dataReg.FindStringSubmatch(rows[i])

		if len(parts) == 0 {
			err := &ErrCannotParseInput{Input: rows[i]}
			t.logg.Warn(err.Error())
			continue
		}

		proto := parts[1]

		reqvQ, err := strconv.Atoi(parts[2])
		if err != nil {
			t.logg.Warn(err.Error())
			continue
		}
		sendQ, err := strconv.Atoi(parts[3])
		if err != nil {
			t.logg.Warn(err.Error())
			continue
		}

		locAddr := parts[4]
		locAddrParts := strings.Split(locAddr, ":")

		var locPort int
		if locAddrParts[len(locAddrParts)-1] != "*" {
			locPort, err = strconv.Atoi(locAddrParts[len(locAddrParts)-1])
			if err != nil {
				t.logg.Warn(err.Error())
				continue
			}
		}

		forAddr := parts[5]
		forAddrParts := strings.Split(forAddr, ":")

		var forPort int
		if forAddrParts[len(forAddrParts)-1] != "*" {
			forPort, err = strconv.Atoi(forAddrParts[len(forAddrParts)-1])
			if err != nil {
				t.logg.Warn(err.Error())
				continue
			}
		}

		state := parts[6]
		pidCmd := parts[7]
		pidCmdParts := strings.SplitN(pidCmd, "/", 2)

		var PID int
		var program string
		if pidCmdParts[0] != "-" {
			PID, err = strconv.Atoi(pidCmdParts[0])
			if err != nil {
				t.logg.Warn(err.Error())
				continue
			}
			program = pidCmdParts[1]
		}

		connects = append(connects, NetStatRow{
			Proto:       proto,
			RecvQ:       reqvQ,
			SendQ:       sendQ,
			LocalAddr:   locAddr,
			LocalPort:   locPort,
			ForeignAddr: forAddr,
			ForeignPort: forPort,
			State:       state,
			PID:         PID,
			Programm:    program,
		})
	}

	return connects
}

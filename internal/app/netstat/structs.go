package appnetstat

import "fmt"

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Parser interface {
	Parse(in string) ([]NetStatRow, error)
}

type NetStatRow struct {
	Proto, LocalAddr, ForeignAddr, State, Programm string
	RecvQ, SendQ, PID, LocalPort, ForeignPort      int
}

type ConnectData struct {
	Infos  []ConnectInfo
	States []ConnectState
}

type ConnectInfo struct {
	ID       string
	Command  string // top
	Pid      int    // 77
	User     string // ?
	Protocol string // TCP
	Port     int    // 40349
}

type ConnectState struct {
	ID       string
	Protocol string // tcp
	State    string // listen
}

type ErrCannotParseInput struct {
	Input string
}

func (e *ErrCannotParseInput) Error() string {
	return fmt.Sprintf("cannot parse %s", e.Input)
}

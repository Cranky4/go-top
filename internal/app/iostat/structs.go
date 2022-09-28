package appiostat

import "fmt"

type DiskIO struct {
	Device    string  // nvme0n1
	Tps, Kbps float32 // 52.86, 665.63 + 780.25
}

type IostatRow struct {
	Device                   string
	Tps, KbpsRead, KbpsWrite float32
	KbRead, KbWrite          int
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

// errors
type ErrCannotParseInput struct {
	Input string
}

func (e *ErrCannotParseInput) Error() string {
	return fmt.Sprintf("cannot parse %s", e.Input)
}

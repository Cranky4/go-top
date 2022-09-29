package apptop

import "fmt"

type CPU struct {
	Avg   CPUAvg
	State CPUState
}

type CPUAvg struct {
	Min     float32 // The first value depicts the average load on the CPU for the last minute.
	Five    float32 // The second gives us the average load for the last 5-minute interval
	Fifteen float32 // The third value gives us the 15-minute average load
}

type CPUState struct {
	User, System, Idle float32
}

type Parser interface {
	Parse(in string) (CPU, error)
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type ErrCannotParseInput struct {
	Input string
}

func (e *ErrCannotParseInput) Error() string {
	return fmt.Sprintf("cannot parse %s", e.Input)
}

package apptop

import "time"

type Cpu struct {
	StartTime, FinishTime time.Time
	Avg                   CpuAvg
	State                 CpuState
}

type CpuAvg struct {
	Min     float32 // The first value depicts the average load on the CPU for the last minute.
	Five    float32 // The second gives us the average load for the last 5-minute interval
	Fifteen float32 // The third value gives us the 15-minute average load
}

type CpuState struct {
	User, System, Idle float32
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

package appiostat

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type TParser struct{}

func (p *TParser) Parse(in string) []IostatRow {
	return []IostatRow{
		{
			Device:    "nvme0n1",
			Tps:       72.84,
			KbpsRead:  329.34,
			KbpsWrite: 2369.18,
			KbRead:    5943837,
			KbWrite:   42758633,
		},
	}
}

type TLogger struct {
	messages []string
}

func (l *TLogger) Info(msg string) {
	l.messages = append(l.messages, msg)
}

func (l *TLogger) Warn(msg string) {
	l.messages = append(l.messages, msg)
}

func (l *TLogger) Error(msg string) {
	l.messages = append(l.messages, msg)
}

func (l *TLogger) Debug(msg string) {
	l.messages = append(l.messages, msg)
}

func TestRun(t *testing.T) {
	logg := new(TLogger)
	parser := new(TParser)

	runner := New("ls", logg, parser)

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	ch := runner.Run(ctx, 1, 1)

	ex := []DiskIO{
		{
			Device: "nvme0n1",
			Tps:    72.84,
			Kbps:   329.34 + 2369.18,
		},
	}

	var i int
	for v := range ch {
		i++
		require.Equal(t, ex, v)

		if i == 2 {
			cancelFn()
		}
	}

	logEx := []string{
		"[IostatRunner] started",
		"[IostatRunner] warmed up",
		"[IostatRunner] collected",
	}

	require.Equal(t, logEx[:3], logg.messages[:3])
}

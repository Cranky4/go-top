package apptop

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type TParser struct{}

func (p *TParser) Parse(in string) (CPU, error) {
	return CPU{
		Avg: CPUAvg{
			Min:     float32(0.23),
			Five:    float32(0.19),
			Fifteen: float32(0.13),
		},
		State: CPUState{
			User:   float32(0.02),
			System: float32(0.12),
			Idle:   float32(0.86),
		},
	}, nil
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

	runner := New("ls", parser, logg)

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	ch := runner.Run(ctx, 1, 1)

	ex := CPU{
		Avg: CPUAvg{
			Min:     float32(0.23),
			Five:    float32(0.19),
			Fifteen: float32(0.13),
		},
		State: CPUState{
			User:   float32(0.02),
			System: float32(0.12),
			Idle:   float32(0.86),
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
		"[TopRunner] started",
		"[TopRunner] warmed up",
		"[TopRunner] collected",
	}

	require.Equal(t, logEx[:3], logg.messages[:3])
}

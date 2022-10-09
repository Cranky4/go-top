package appnetstat

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type TParser struct{}

func (p *TParser) Parse(in string) []NetStatRow {
	return []NetStatRow{
		{
			Proto:       "tcp",
			RecvQ:       0,
			SendQ:       0,
			LocalAddr:   "127.0.0.11:36269",
			LocalPort:   36269,
			ForeignAddr: "0.0.0.0:*",
			ForeignPort: 0,
			State:       "LISTEN",
			PID:         88,
			Programm:    "cmd",
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

	ex := ConnectData{
		Infos: []ConnectInfo{
			{
				ID:       "127.0.0.11:362690.0.0.0:*tcp",
				Command:  "cmd",
				Pid:      88,
				User:     "",
				Protocol: "tcp",
				Port:     36269,
			},
		},
		States: []ConnectState{
			{
				ID:       "127.0.0.11:362690.0.0.0:*tcp",
				Protocol: "tcp",
				State:    "LISTEN",
			},
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
		"[NetstatRunner] started",
		"[NetstatRunner] warmed up",
		"[NetstatRunner] collected",
	}

	require.Equal(t, logEx[:3], logg.messages[:3])
}

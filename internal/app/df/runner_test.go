package appdf

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type TParser struct{}

func (p *TParser) ParseBytes(in string) []DiskInfo {
	return []DiskInfo{
		{
			Name:           "/dev/nvme0n1p8",
			AvailableBytes: 11318152,
			UsedBytes:      48946056,
			UsageBytes:     82,
		},
	}
}

func (p *TParser) ParseInodes(in string) []DiskInfo {
	return []DiskInfo{
		{
			Name:            "/dev/nvme0n1p8",
			AvailableInodes: 0,
			UsedInodes:      0,
			UsageInodes:     0,
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

	ex := []DiskInfo{
		{
			Name:            "/dev/nvme0n1p8",
			AvailableBytes:  11318152,
			UsedBytes:       48946056,
			UsageBytes:      82,
			AvailableInodes: 0,
			UsedInodes:      0,
			UsageInodes:     0,
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
		"[DfRunner] started",
		"[DfRunner] warmed up",
		"[DfRunner] collected",
	}

	require.Equal(t, logEx[:3], logg.messages[:3])
}

package apptcpdump

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type TParser struct{}

func (p *TParser) Parse(in string) ([]TCPDumpLine, error) {
	return []TCPDumpLine{
		{
			Time:        time.Now(),
			Type:        "IP6",
			Source:      "fe80::42:c1ff:fe1d:f67d.5353",
			Destination: "ff02::fb.5353",
			Protocol:    "UDP",
			Bytes:       118,
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

	runner := New("ls", "ls", logg, parser)

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	ch := runner.Run(ctx, 1, 1)

	ex := TopTalkers{
		ByProtocol: []TopTalkerByProtocol{
			{
				Protocol: "UDP",
				Bytes:    118,
				Percent:  "100%",
			},
		},
		ByTraffic: []TopTalkerByTraffic{
			{
				Source:         "fe80::42:c1ff:fe1d:f67d.5353",
				Destination:    "ff02::fb.5353",
				Protocol:       "UDP",
				BytesPerSecond: 118,
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
}

package apptcpdump

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"sort"
	"time"
)

type TCPDumpRunner struct {
	commandPath, timeoutPath string
	parser                   Parser
	logg                     Logger
}

func New(timeoutPath, commandPath string, logg Logger, parser Parser) *TCPDumpRunner {
	return &TCPDumpRunner{
		timeoutPath: timeoutPath,
		commandPath: commandPath,
		parser:      parser,
		logg:        logg,
	}
}

func (t *TCPDumpRunner) Run(ctx context.Context, m, n int) chan TopTalkers {
	ch := make(chan TopTalkers)
	t.logg.Debug("[TcpDumpRunner] started")
	started := time.Now()

	go func() {
		defer close(ch)
		data := make([]TCPDumpLine, 0, m)

		data, err := t.collect(ctx, m, data)
		if err != nil {
			t.logg.Error(
				fmt.Sprintf("[TcpDumpRunner] err: %s", err),
			)
			return
		}

		// warming up
		talkers := t.calculate(data)

		select {
		case <-ctx.Done():
			return
		case ch <- talkers:
			t.logg.Debug("[TcpDumpRunner] collected")
		}

		// collect
		for {
			dur, err := time.ParseDuration(fmt.Sprintf("%ds", n))
			if err != nil {
				t.logg.Error(err.Error())
				return
			}
			started = started.Add(dur)
			data = t.cleanOldLines(data, started)

			data, err = t.collect(ctx, n, data)
			if err != nil {
				t.logg.Error(
					fmt.Sprintf("[TcpDumpRunner] err: %s", err),
				)
				return
			}

			talkers := t.calculate(data)

			select {
			case <-ctx.Done():
				return
			case ch <- talkers:
				t.logg.Debug("[TcpDumpRunner] collected")
			}
		}
	}()

	return ch
}

func (t *TCPDumpRunner) cleanOldLines(data []TCPDumpLine, threshold time.Time) []TCPDumpLine {
	for i, l := range data {
		if l.Time.Before(threshold) {
			continue
		}

		return data[i:]
	}

	return make([]TCPDumpLine, 0, len(data))
}

func (t *TCPDumpRunner) collect(ctx context.Context, seconds int, data []TCPDumpLine) ([]TCPDumpLine, error) {
	cmd := exec.CommandContext( //nolint:gosec
		ctx,
		t.timeoutPath,
		"--preserve-status",
		fmt.Sprintf("%ds", seconds),
		t.commandPath,
		"-ntq", "-i", "any", "-Q", "inout", "-ttt", "-l",
	)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	lines, err := t.parser.Parse(out.String())
	if err != nil {
		return nil, err
	}

	return append(data, lines...), nil
}

func (t *TCPDumpRunner) calculate(lines []TCPDumpLine) TopTalkers {
	byProtocolMap := make(map[string]int) // [protocol]bytes

	byTrafficSecondsMap := make(map[string]map[string]struct{}) // [pseudoId][secondsCount]struct
	byTrafficBytesMap := make(map[string]int)                   // [pseudoId]totalBytes
	byTrafficLineMap := make(map[string]TCPDumpLine)            // [pseudoId]line

	var totalBytes int
	for _, l := range lines {
		_, ex := byProtocolMap[l.Protocol]
		if !ex {
			byProtocolMap[l.Protocol] = 0
		}

		byProtocolMap[l.Protocol] += l.Bytes
		totalBytes += l.Bytes

		id := l.Source + l.Destination + l.Protocol
		if _, ex = byTrafficLineMap[id]; !ex {
			byTrafficLineMap[id] = l
		}

		if _, ex = byTrafficBytesMap[id]; ex {
			byTrafficBytesMap[id] += l.Bytes
		} else {
			byTrafficBytesMap[id] = 0
		}

		time := l.Time.Format("2006-01-02 15:04:05")
		if _, ex = byTrafficSecondsMap[id]; ex {
			if _, ex = byTrafficSecondsMap[id][time]; !ex {
				byTrafficSecondsMap[id][time] = struct{}{}
			}
		} else {
			byTrafficSecondsMap[id] = make(map[string]struct{})
			byTrafficSecondsMap[id][time] = struct{}{}
		}
	}

	byProtocol := make([]TopTalkerByProtocol, 0, len(byProtocolMap))
	for protocol, bytes := range byProtocolMap {
		byProtocol = append(byProtocol, TopTalkerByProtocol{
			Protocol: protocol,
			Bytes:    bytes,
			Percent:  fmt.Sprintf("%.2f%%", float32(bytes)/float32(totalBytes)*float32(100)),
		})
	}

	sort.Slice(byProtocol, func(i, j int) bool {
		return byProtocol[i].Percent > byProtocol[j].Percent
	})

	byTraffic := make([]TopTalkerByTraffic, 0, len(byTrafficSecondsMap))
	for id, bytes := range byTrafficBytesMap {
		if bytes == 0 {
			continue
		}
		bps := float32(bytes) / float32(len(byTrafficSecondsMap[id]))

		line := byTrafficLineMap[id]
		byTraffic = append(byTraffic, TopTalkerByTraffic{
			Source:         line.Source,
			Destination:    line.Destination,
			Protocol:       line.Protocol,
			BytesPerSecond: bps,
		})
	}

	sort.Slice(byTraffic, func(i, j int) bool {
		return byTraffic[i].BytesPerSecond > byTraffic[j].BytesPerSecond
	})

	return TopTalkers{
		ByProtocol: byProtocol,
		ByTraffic:  byTraffic,
	}
}

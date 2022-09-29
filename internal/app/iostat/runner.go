package appiostat

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

type IostatRunner struct {
	commandPath string
	args        []string
	parser      Parser
	logg        Logger
}

func New(commandPath string, logg Logger, parser Parser) *IostatRunner {
	return &IostatRunner{
		commandPath: commandPath,
		args:        []string{"-d", "-k"},
		parser:      parser,
		logg:        logg,
	}
}

func (t *IostatRunner) Run(ctx context.Context, m, n int) chan []DiskIO {
	ch := make(chan []DiskIO)
	t.logg.Debug("[IostatRunner] started")

	go func() {
		defer close(ch)
		data := make([][]DiskIO, 0, m)

		err := t.collect(ctx, m, &data)
		if err != nil {
			t.logg.Error(
				fmt.Sprintf("[IostatRunner] err: %s", err),
			)
			return
		}

		// warming up
		avg := t.calculateAvg(data)

		select {
		case <-ctx.Done():
			return
		case ch <- avg:
			t.logg.Debug("[IostatRunner] warmed up")
			data = data[n:]
		}

		// collect
		for {
			err = t.collect(ctx, n, &data)
			if err != nil {
				t.logg.Error(
					fmt.Sprintf("[IostatRunner] err: %s", err),
				)
				return
			}

			avg := t.calculateAvg(data)

			select {
			case <-ctx.Done():
				return
			case ch <- avg:
				t.logg.Debug("[IostatRunner] collected")
				data = data[n:]
			}
		}
	}()

	return ch
}

func (t *IostatRunner) collect(ctx context.Context, seconds int, data *[][]DiskIO) error {
	for i := 0; i < seconds; i++ {
		select {
		case <-ctx.Done():
		default:
			cmd := exec.CommandContext(ctx, t.commandPath, t.args...) //nolint:gosec

			var out bytes.Buffer
			cmd.Stdout = &out

			if err := cmd.Run(); err != nil {
				return err
			}

			rows, err := t.parser.Parse(out.String())
			if err != nil {
				return err
			}

			ios := make([]DiskIO, 0, len(rows))
			for _, r := range rows {
				ios = append(ios, DiskIO{
					Device: r.Device,
					Tps:    r.Tps,
					Kbps:   r.KbpsRead + r.KbpsWrite,
				})
			}

			*data = append(*data, ios)

			if seconds > 1 {
				time.Sleep(1 * time.Second)
			}
		}
	}

	return nil
}

func (t *IostatRunner) calculateAvg(disks [][]DiskIO) []DiskIO {
	devices := make(map[string][][]float32)

	for _, dd := range disks {
		for _, d := range dd {
			_, ex := devices[d.Device]
			if !ex {
				devices[d.Device] = make([][]float32, 2)
			}

			devices[d.Device][0] = append(devices[d.Device][0], d.Tps)
			devices[d.Device][1] = append(devices[d.Device][1], d.Kbps)
		}
	}

	result := make([]DiskIO, 0, len(devices))

	for d, values := range devices {
		tpss := values[0]
		var tpsSum float32 = 0
		for _, tps := range tpss {
			tpsSum += tps
		}

		kbpss := values[1]
		var kbpsSum float32 = 0
		for _, kbps := range kbpss {
			kbpsSum += kbps
		}

		result = append(result, DiskIO{
			Device: d,
			Tps:    tpsSum / float32(len(values[0])),
			Kbps:   kbpsSum / float32(len(values[1])),
		})
	}

	return result
}

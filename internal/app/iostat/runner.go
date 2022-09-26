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
	parser      *IostatParser
	logg        Logger
}

func New(commandPath string, logg Logger) *IostatRunner {
	return &IostatRunner{
		commandPath: commandPath,
		args:        []string{"-d", "-k"},
		parser:      NewParser(),
		logg:        logg,
	}
}

func (t *IostatRunner) Run(ctx context.Context, M, N uint32) chan []DiskIO {
	ch := make(chan []DiskIO)
	t.logg.Debug("[IostatRunner] started")

	go func() {
		defer close(ch)
		data := make([][]DiskIO, 0, M)

		err := t.collect(ctx, M, &data)
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
			data = data[N:]
		}

		// collect
		for {
			err = t.collect(ctx, N, &data)
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
				data = data[N:]
			}
		}
	}()

	return ch
}

func (t *IostatRunner) collect(ctx context.Context, seconds uint32, data *[][]DiskIO) error {
	var i uint32

	for i = 0; i < seconds; i++ {
		select {
		case <-ctx.Done():
		default:
			cmd := exec.CommandContext(
				ctx,
				t.commandPath,
				t.args...,
			)

			var out bytes.Buffer
			cmd.Stdout = &out

			if err := cmd.Run(); err != nil {
				return err
			}

			item, err := t.parser.Parse(out.String())
			if err != nil {
				return err
			}

			*data = append(*data, item)

			time.Sleep(1 * time.Second)
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
				devices[d.Device] = make([][]float32, 2, 2)
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

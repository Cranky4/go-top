package apptop

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

type TopRunner struct {
	commandPath string
	args        []string
	parser      Parser
	logg        Logger
}

func New(commandPath string, parser Parser, logg Logger) *TopRunner {
	return &TopRunner{
		commandPath: commandPath,
		args:        []string{"-b", "-n1"},
		parser:      parser,
		logg:        logg,
	}
}

func (t *TopRunner) Run(ctx context.Context, m, n int) chan CPU {
	ch := make(chan CPU)
	t.logg.Debug("[TopRunner] started")

	go func() {
		defer close(ch)
		cpus := make([]CPU, 0, m)

		err := t.collect(ctx, m, &cpus)
		if err != nil {
			t.logg.Error(
				fmt.Sprintf("[TopRunner] err: %s", err),
			)
			return
		}

		// warming up
		avg := t.calculateAvg(cpus)

		select {
		case <-ctx.Done():
			return
		case ch <- avg:
			t.logg.Debug("[TopRunner] warmed up")
			cpus = cpus[n:]
		}

		// collect
		for {
			err = t.collect(ctx, n, &cpus)
			if err != nil {
				t.logg.Error(
					fmt.Sprintf("[TopRunner] err: %s", err),
				)
				return
			}

			avg := t.calculateAvg(cpus)

			select {
			case <-ctx.Done():
				return
			case ch <- avg:
				t.logg.Debug("[TopRunner] collected")
				cpus = cpus[n:]
			}
		}
	}()

	return ch
}

func (t *TopRunner) collect(ctx context.Context, seconds int, cpus *[]CPU) error {
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

			cpu, err := t.parser.Parse(out.String())
			if err != nil {
				return err
			}

			*cpus = append(*cpus, cpu)

			if seconds > 1 {
				time.Sleep(1 * time.Second)
			}
		}
	}

	return nil
}

func (t *TopRunner) calculateAvg(cpus []CPU) CPU {
	var sumAvgMin, sumAvgFive, sumAvgFifteen, sumStateUser, sumStateSystem, sumStateIdle float32
	for _, c := range cpus {
		sumAvgMin += c.Avg.Min
		sumAvgFive += c.Avg.Five
		sumAvgFifteen += c.Avg.Fifteen
		sumStateUser += c.State.User
		sumStateSystem += c.State.System
		sumStateIdle += c.State.Idle
	}

	devizor := float32(len(cpus))

	return CPU{
		Avg: CPUAvg{
			Min:     sumAvgMin / devizor,
			Five:    sumAvgFive / devizor,
			Fifteen: sumAvgFifteen / devizor,
		},
		State: CPUState{
			User:   sumStateUser / devizor,
			System: sumStateSystem / devizor,
			Idle:   sumStateIdle / devizor,
		},
	}
}

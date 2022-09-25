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
	parser      *TopParser
	logg        Logger
}

func New(commandPath string, logg Logger) *TopRunner {
	return &TopRunner{
		commandPath: commandPath,
		args:        []string{"-b", "-n1"},
		parser:      NewParser(),
		logg:        logg,
	}
}

func (t *TopRunner) Run(ctx context.Context, M, N uint32, StartTime time.Time) chan Cpu {
	ch := make(chan Cpu)
	t.logg.Debug("[TopRunner] started")

	go func(startTime time.Time) {
		defer close(ch)
		cpus := make([]Cpu, 0, M)

		err := t.collect(ctx, M, &cpus)
		if err != nil {
			t.logg.Error(
				fmt.Sprintf("[TopRunner] err: %s", err),
			)
			return
		}

		// warming up
		avgCpu := t.calculateAvg(cpus)
		avgCpu.StartTime = startTime
		d, _ := time.ParseDuration(fmt.Sprintf("%ds", M))
		avgCpu.FinishTime = startTime.Add(d)
		startTime = avgCpu.FinishTime

		select {
		case <-ctx.Done():
			return
		case ch <- avgCpu:
			t.logg.Debug("[TopRunner] warmed up")
			cpus = cpus[N:]
		}

		// collect
		for {
			err = t.collect(ctx, N, &cpus)
			if err != nil {
				t.logg.Error(
					fmt.Sprintf("[TopRunner] err: %s", err),
				)
				return
			}

			avgCpu := t.calculateAvg(cpus)
			avgCpu.StartTime = startTime
			d, _ := time.ParseDuration(fmt.Sprintf("%ds", N))
			avgCpu.FinishTime = startTime.Add(d)

			select {
			case <-ctx.Done():
				return
			case ch <- avgCpu:
				t.logg.Debug("[TopRunner] collected")
				cpus = cpus[N:]
				startTime = avgCpu.FinishTime
			}
		}
	}(StartTime)

	return ch
}

func (t *TopRunner) collect(ctx context.Context, seconds uint32, cpus *[]Cpu) error {
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

			cpu, err := t.parser.Parse(out.String())
			if err != nil {
				return err
			}

			*cpus = append(*cpus, cpu)

			time.Sleep(1 * time.Second)
		}
	}

	return nil
}

func (t *TopRunner) calculateAvg(cpus []Cpu) Cpu {
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

	return Cpu{
		Avg: CpuAvg{
			Min:     sumAvgMin / devizor,
			Five:    sumAvgFive / devizor,
			Fifteen: sumAvgFifteen / devizor,
		},
		State: CpuState{
			User:   sumStateUser / devizor,
			System: sumStateSystem / devizor,
			Idle:   sumStateIdle / devizor,
		},
	}
}
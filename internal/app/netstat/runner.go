package appnetstat

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

type NetstatRunner struct {
	commandPath string
	parser      Parser
	logg        Logger
}

func New(commandPath string, logg Logger, parser Parser) *NetstatRunner {
	return &NetstatRunner{
		commandPath: commandPath,
		parser:      parser,
		logg:        logg,
	}
}

func (t *NetstatRunner) Run(ctx context.Context, M, N uint32) chan ConnectData {
	ch := make(chan ConnectData)
	t.logg.Debug("[NetstatRunner] started")

	go func() {
		defer close(ch)
		infos := make([][]ConnectInfo, 0, M)
		states := make([][]ConnectState, 0, M)

		err := t.collect(ctx, M, &infos, &states)
		if err != nil {
			t.logg.Error(
				fmt.Sprintf("[NetstatRunner] err: %s", err),
			)
			return
		}

		// warming up
		connectData := ConnectData{
			Infos:  t.calculateUniques(infos),
			States: t.calculateStates(states),
		}

		select {
		case <-ctx.Done():
			return
		case ch <- connectData:
			t.logg.Debug("[NetstatRunner] warmed up")
			infos = infos[N:]
		}

		// collect
		for {
			err = t.collect(ctx, N, &infos, &states)
			if err != nil {
				t.logg.Error(
					fmt.Sprintf("[NetstatRunner] err: %s", err),
				)
				return
			}

			connectData := ConnectData{
				Infos:  t.calculateUniques(infos),
				States: t.calculateStates(states),
			}

			select {
			case <-ctx.Done():
				return
			case ch <- connectData:
				t.logg.Debug("[NetstatRunner] collected")
				infos = infos[N:]
				states = states[N:]
			}
		}
	}()

	return ch
}

func (t *NetstatRunner) collect(
	ctx context.Context,
	seconds uint32,
	connectsInfos *[][]ConnectInfo,
	connectsStates *[][]ConnectState,
) error {
	var i uint32

	for i = 0; i < seconds; i++ {
		select {
		case <-ctx.Done():
		default:
			cmd := exec.CommandContext(
				ctx,
				t.commandPath,
				"-lntap",
			)

			var out bytes.Buffer
			cmd.Stdout = &out

			if err := cmd.Run(); err != nil {
				return err
			}

			netStatRows, err := t.parser.Parse(out.String())
			if err != nil {
				return err
			}

			connects := make([]ConnectInfo, 0, len(netStatRows))
			states := make([]ConnectState, 0, len(netStatRows))
			for _, r := range netStatRows {
				connects = append(connects, ConnectInfo{
					ID:       r.LocalAddr + r.ForeignAddr + r.Proto,
					Command:  r.Programm,
					Pid:      r.PID,
					Protocol: r.Proto,
					Port:     r.LocalPort,
				})

				states = append(states, ConnectState{
					ID:       r.LocalAddr + r.ForeignAddr + r.Proto,
					Protocol: r.Proto,
					State:    r.State,
				})
			}

			*connectsInfos = append(*connectsInfos, connects)
			*connectsStates = append(*connectsStates, states)

			time.Sleep(1 * time.Second)
		}
	}

	return nil
}

func (t *NetstatRunner) calculateUniques(connsSets [][]ConnectInfo) []ConnectInfo {
	connections := make(map[string]ConnectInfo)

	// optimistical find uniques connections
	for _, conns := range connsSets {
		for _, conn := range conns {
			_, ex := connections[conn.ID]
			if !ex {
				connections[conn.ID] = conn
			}
		}
	}

	result := make([]ConnectInfo, 0, len(connections))

	for _, conn := range connections {
		result = append(result, conn)
	}

	return result
}

func (t *NetstatRunner) calculateStates(connsSets [][]ConnectState) []ConnectState {
	connections := make(map[string]ConnectState)

	// optimistical find uniques
	for _, conns := range connsSets {
		for _, conn := range conns {
			_, ex := connections[conn.ID]
			if !ex {
				connections[conn.ID] = conn
			}
		}
	}

	result := make([]ConnectState, 0, len(connections))

	for _, conn := range connections {
		result = append(result, conn)
	}

	return result
}

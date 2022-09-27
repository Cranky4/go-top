package apptcpdump

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

type TcpDumpRunner struct {
	commandPath string
	args        []string
	parser      Parser
	logg        Logger
}

func New(commandPath string, logg Logger, parser Parser) *TcpDumpRunner {
	return &TcpDumpRunner{
		commandPath: commandPath,
		args:        []string{"-ntq", "-i", "any", "-Q", "inout", "-ttt", "-l"},
		parser:      parser,
		logg:        logg,
	}
}

func (t *TcpDumpRunner) Run(ctx context.Context, M, N uint32) chan TopTalkers {
	ch := make(chan TopTalkers)
	t.logg.Debug("[TcpDumpRunner] started")

	go func() {
		defer close(ch)
		data := make([]TcpDumpLine, 0, M)

		err := t.collect(ctx, M, &data)
		if err != nil {
			t.logg.Error(
				fmt.Sprintf("[TcpDumpRunner] err: %s", err),
			)
			return
		}

		t.logg.Debug(fmt.Sprintf("%#v", data))

		// warming up
		// avg := t.calculateAvg(data)

		// select {
		// case <-ctx.Done():
		// 	return
		// case ch <- avg:
		// 	t.logg.Debug("[TcpDumpRunner] warmed up")
		// 	data = data[N:]
		// }

		// // collect
		// for {
		// 	err = t.collect(ctx, N, &data)
		// 	if err != nil {
		// 		t.logg.Error(
		// 			fmt.Sprintf("[TcpDumpRunner] err: %s", err),
		// 		)
		// 		return
		// 	}

		// 	avg := t.calculateAvg(data)

		// 	select {
		// 	case <-ctx.Done():
		// 		return
		// 	case ch <- avg:
		// 		t.logg.Debug("[TcpDumpRunner] collected")
		// 		data = data[N:]
		// 	}
		// }
	}()

	return ch
}

func (t *TcpDumpRunner) collect(ctx context.Context, seconds uint32, data *[]TcpDumpLine) error {
	dur, err := time.ParseDuration(fmt.Sprintf("%ds", seconds))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, dur)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		t.commandPath,
		t.args...,
	)

	cmdPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Run(); err != nil {
		oe, _ := cmd.Output()
		t.logg.Debug(string(oe))
		return err
	}

	buffer := make([]byte, 4096, 4096)
	readed, err := cmdPipe.Read(buffer)
	if err != nil {
		return err
	}

	t.logg.Debug(string(buffer[0:readed]))

	// item, err := t.parser.Parse(out.String())
	// if err != nil {
	// 	return err
	// }

	// *data = append(*data, item)

	// time.Sleep(1 * time.Second)

	return nil
}

// func (t *TcpDumpRunner) calculateAvg(lines []TcpDumpLine) []TcpDumpLine {
// 	talkersBytes := make(map[string][]TcpDumpLine)

// 	for _, l := range lines {
// 		_, ex := talkersBytes[l.Protocol]
// 		if !ex {
// 			talkersBytes[l.Protocol] = make([]TcpDumpLine, 1, 1)
// 		}

// 		talkersBytes[l.Protocol] = append(talkersBytes[l.Protocol], l)
// 	}

// 	result := make([]TcpDumpLine, 0, len(talkersBytes))

// 	for d, values := range talkers {
// 		tpss := values[0]
// 		var tpsSum float32 = 0
// 		for _, tps := range tpss {
// 			tpsSum += tps
// 		}

// 		kbpss := values[1]
// 		var kbpsSum float32 = 0
// 		for _, kbps := range kbpss {
// 			kbpsSum += kbps
// 		}

// 		result = append(result, DiskIO{
// 			Device: d,
// 			Tps:    tpsSum / float32(len(values[0])),
// 			Kbps:   kbpsSum / float32(len(values[1])),
// 		})
// 	}

// 	return result
// }

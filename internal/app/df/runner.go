package appdf

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

type IostatRunner struct {
	commandPath string
	parser      Parser
	logg        Logger
}

func New(commandPath string, logg Logger, parser Parser) *IostatRunner {
	return &IostatRunner{
		commandPath: commandPath,
		parser:      parser,
		logg:        logg,
	}
}

func (t *IostatRunner) Run(ctx context.Context, warmingUpTime, shapshotPeriod int) chan []DiskInfo {
	ch := make(chan []DiskInfo)
	t.logg.Debug("[DfRunner] started")

	go func() {
		defer close(ch)
		data := make([][]DiskInfo, 0, warmingUpTime)

		err := t.collect(ctx, warmingUpTime, &data)
		if err != nil {
			t.logg.Error(
				fmt.Sprintf("[DfRunner] err: %s", err),
			)
			return
		}

		// warming up
		avg := t.calculateAvg(data)

		select {
		case <-ctx.Done():
			return
		case ch <- avg:
			t.logg.Debug("[DfRunner] warmed up")
			data = data[shapshotPeriod:]
		}

		// collect
		for {
			err = t.collect(ctx, shapshotPeriod, &data)
			if err != nil {
				t.logg.Error(
					fmt.Sprintf("[DfRunner] err: %s", err),
				)
				return
			}

			avg := t.calculateAvg(data)

			select {
			case <-ctx.Done():
				return
			case ch <- avg:
				t.logg.Debug("[DfRunner] collected")
				data = data[shapshotPeriod:]
			}
		}
	}()

	return ch
}

func (t *IostatRunner) collect(ctx context.Context, seconds int, disks *[][]DiskInfo) error {
	for i := 0; i < seconds; i++ {
		select {
		case <-ctx.Done():
		default:
			// df -k
			cmd := exec.CommandContext(ctx, t.commandPath, "-k") //nolint:gosec

			var out bytes.Buffer
			cmd.Stdout = &out

			if err := cmd.Run(); err != nil {
				return err
			}

			diskBytes := t.parser.ParseBytes(out.String())

			// df -i
			cmd = exec.CommandContext(ctx, t.commandPath, "-i") //nolint:gosec
			cmd.Stdout = &out

			if err := cmd.Run(); err != nil {
				return err
			}

			diskInodes := t.parser.ParseInodes(out.String())

			// merge output
			merged := t.merge(diskBytes, diskInodes)

			currentDiscs := make([]DiskInfo, 0, len(merged))
			for _, di := range merged {
				currentDiscs = append(currentDiscs, di)
			}
			*disks = append(*disks, currentDiscs)

			if seconds > 1 {
				time.Sleep(1 * time.Second)
			}
		}
	}

	return nil
}

func (*IostatRunner) merge(diskBytes []DiskInfo, diskInodes []DiskInfo) map[string]DiskInfo {
	merged := make(map[string]DiskInfo)
	for _, d := range diskBytes {
		merged[d.Name] = DiskInfo{
			Name:           d.Name,
			AvailableBytes: d.AvailableBytes,
			UsedBytes:      d.UsedBytes,
			UsageBytes:     d.UsageBytes,
		}
	}

	for _, v := range diskInodes {
		d, ex := merged[v.Name]
		if ex {
			d.AvailableInodes = v.AvailableInodes
			d.UsedInodes = v.UsedInodes
			d.UsageInodes = v.UsageInodes

			merged[v.Name] = d
		}
	}

	return merged
}

func (t *IostatRunner) calculateAvg(disks [][]DiskInfo) []DiskInfo {
	devices := make(map[string][][]int)

	for _, dd := range disks {
		for _, d := range dd {
			_, ex := devices[d.Name]
			if !ex {
				devices[d.Name] = make([][]int, 6)
			}

			devices[d.Name][0] = append(devices[d.Name][0], d.AvailableBytes)
			devices[d.Name][1] = append(devices[d.Name][1], d.UsedBytes)
			devices[d.Name][2] = append(devices[d.Name][2], d.UsageBytes)
			devices[d.Name][3] = append(devices[d.Name][3], d.AvailableInodes)
			devices[d.Name][4] = append(devices[d.Name][4], d.UsedInodes)
			devices[d.Name][5] = append(devices[d.Name][5], d.UsageInodes)
		}
	}

	result := make([]DiskInfo, 0, len(devices))

	for d, values := range devices {
		bytesAvailable := values[0]
		var bytesAvailableSum int
		for _, v := range bytesAvailable {
			bytesAvailableSum += v
		}

		usedBytes := values[1]
		var usedBytesSum int
		for _, v := range usedBytes {
			usedBytesSum += v
		}

		usageBytes := values[2]
		var usageBytesSum int
		for _, v := range usageBytes {
			usageBytesSum += v
		}

		inodeAvailable := values[3]
		var inodeAvailableSum int
		for _, v := range inodeAvailable {
			inodeAvailableSum += v
		}

		usedInodes := values[4]
		var usedInodesSum int
		for _, v := range usedInodes {
			usedInodesSum += v
		}

		usageInodes := values[5]
		var usageInodesSum int
		for _, v := range usageInodes {
			usageInodesSum += v
		}

		result = append(result, DiskInfo{
			Name:            d,
			AvailableBytes:  bytesAvailableSum / len(bytesAvailable),
			UsedBytes:       usedBytesSum / len(usedBytes),
			UsageBytes:      usageBytesSum / len(usageBytes),
			AvailableInodes: inodeAvailableSum / len(inodeAvailable),
			UsedInodes:      usedInodesSum / len(usedInodes),
			UsageInodes:     usageInodesSum / len(usageInodes),
		})
	}

	return result
}

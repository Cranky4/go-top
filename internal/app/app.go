package app

import (
	"context"
	"fmt"

	appdf "github.com/Cranky4/go-top/internal/app/df"
	appiostat "github.com/Cranky4/go-top/internal/app/iostat"
	appnetstat "github.com/Cranky4/go-top/internal/app/netstat"
	apptop "github.com/Cranky4/go-top/internal/app/top"
)

type App struct {
	ctx  context.Context
	conf Config
	logg Logger
	snps map[int64]Snapshot
}

func New(ctx context.Context, conf Config, logg Logger) *App {
	logg.Debug(fmt.Sprintf("[Config] %#v", conf))
	return &App{ctx: ctx, conf: conf, logg: logg}
}

func (t *App) Start(warmUpTime, recordPeriod uint32) <-chan Snapshot {
	ch := make(chan Snapshot)

	go t.work(warmUpTime, recordPeriod, ch)

	return ch
}

func (t *App) work(M, N uint32, ch chan Snapshot) {
	defer close(ch)
	t.logg.Debug("[APP] started... Hello!")

	var cpuCh chan apptop.Cpu
	if t.conf.Metrics.Cpu {
		topRunner := apptop.New(t.conf.App.TopPath, t.logg)
		cpuCh = topRunner.Run(t.ctx, M, N)
	} else {
		t.logg.Debug("[APP] CPU metric is disabled")
	}

	var disksIOCh chan []appiostat.DiskIO
	var disksStatCh chan []appdf.DiskInfo
	if t.conf.Metrics.Disks {
		iostatRunner := appiostat.New(t.conf.App.IostatPath, t.logg)
		disksIOCh = iostatRunner.Run(t.ctx, M, N)

		dfRunner := appdf.New(t.conf.App.DfPath, t.logg, appdf.NewParser(t.logg))
		disksStatCh = dfRunner.Run(t.ctx, M, N)
	} else {
		t.logg.Debug("[APP] discs metric is disabled")
	}

	// chan
	if t.conf.Metrics.Network {
		// tcpDumpRunner := apptcpdump.New(t.conf.App.TcpDumpPath, t.logg, apptcpdump.NewParser())
		// chan = tcpDumpRunner.Run(t.ctx, M, N)
	} else {
		t.logg.Debug("[APP] network metric is disabled")
	}

	var connsCh chan appnetstat.ConnectData
	if t.conf.Metrics.Connections {
		netStatRunner := appnetstat.New(t.conf.App.NetStatPath, t.logg, appnetstat.NewParser(t.logg))
		connsCh = netStatRunner.Run(t.ctx, M, N)
	} else {
		t.logg.Debug("[APP] connection metrics disabled")
	}

L:
	for {
		select {
		case <-t.ctx.Done():
			break L
		default:
			var cpu apptop.Cpu
			if t.conf.Metrics.Cpu {
				t.logg.Debug("[APP] waiting cpu")
				cpu = <-cpuCh
			}

			var disksIO []appiostat.DiskIO
			var disksInfo []appdf.DiskInfo
			if t.conf.Metrics.Disks {
				t.logg.Debug("[APP] waiting disks io")
				disksIO = <-disksIOCh

				t.logg.Debug("[APP] waiting disks stats")
				disksInfo = <-disksStatCh
			}

			var conns appnetstat.ConnectData
			if t.conf.Metrics.Connections {
				t.logg.Debug("[APP] waiting connections")
				conns = <-connsCh
			}

			ch <- Snapshot{
				Cpu:            cpu,
				DisksIO:        disksIO,
				DisksInfo:      disksInfo,
				ConnectsInfo:   conns.Infos,
				ConnectsStates: conns.States,
			}
		}
	}

	t.logg.Info("[APP] finished... Good bye!")
}

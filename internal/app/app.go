package app

import (
	"context"
	"time"

	appdf "github.com/Cranky4/go-top/internal/app/df"
	appiostat "github.com/Cranky4/go-top/internal/app/iostat"
	apptop "github.com/Cranky4/go-top/internal/app/top"
)

type App struct {
	ctx  context.Context
	conf Config
	logg Logger
	snps map[int64]Snapshot
}

func New(ctx context.Context, conf Config, logg Logger) *App {
	return &App{ctx: ctx, conf: conf, logg: logg}
}

func (t *App) Start(warmUpTime, recordPeriod uint32) <-chan Snapshot {
	ch := make(chan Snapshot)

	go t.work(warmUpTime, recordPeriod, ch)

	return ch
}

func (t *App) work(M, N uint32, ch chan Snapshot) {
	defer close(ch)
	t.logg.Info("[APP] top collection started...")

	start := time.Now()

	topRunner := apptop.New(t.conf.App.TopPath, t.logg)
	cpuCh := topRunner.Run(t.ctx, M, N)

	iostatRunner := appiostat.New(t.conf.App.IostatPath, t.logg)
	disksIOCh := iostatRunner.Run(t.ctx, M, N)

	dfRunner := appdf.New(t.conf.App.DfPath, t.logg, appdf.NewParser(t.logg))
	disksStatCh := dfRunner.Run(t.ctx, M, N)

L:
	for {
		select {
		case <-t.ctx.Done():
			break L
		default:
			t.logg.Debug("[APP] waiting cpu")
			cpu := <-cpuCh

			t.logg.Debug("[APP] waiting disks io")
			disksIO := <-disksIOCh

			t.logg.Debug("[APP] waiting disks io")
			disksInfo := <-disksStatCh

			ch <- Snapshot{
				StartTime:  start,
				FinishTime: start,
				Cpu:        cpu,
				DisksIO:    disksIO,
				DisksInfo:  disksInfo,
			}
		}
	}

	t.logg.Info("[APP] top collection finished...")
}

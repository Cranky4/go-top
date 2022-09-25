package app

import (
	"context"
	"time"

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
	t.logg.Info("top collection started...")

	topRunner := apptop.New(t.conf.Top.TopPath, t.logg)

	cpuCh := topRunner.Run(t.ctx, M, N, time.Now())

	for cpu := range cpuCh {
		ch <- Snapshot{
			StartTime:  cpu.StartTime,
			FinishTime: cpu.FinishTime,
			Cpu:        cpu,
		}
	}

	t.logg.Info("top collection finished...")
}

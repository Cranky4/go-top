package app

import (
	"context"
	"fmt"

	appdf "github.com/Cranky4/go-top/internal/app/df"
	appiostat "github.com/Cranky4/go-top/internal/app/iostat"
	appnetstat "github.com/Cranky4/go-top/internal/app/netstat"
	apptcpdump "github.com/Cranky4/go-top/internal/app/tcpdump"
	apptop "github.com/Cranky4/go-top/internal/app/top"
)

type App struct {
	ctx  context.Context
	conf Config
	logg Logger
}

func New(ctx context.Context, conf Config, logg Logger) *App {
	logg.Debug(fmt.Sprintf("[Config] %#v", conf))
	return &App{ctx: ctx, conf: conf, logg: logg}
}

func (t *App) Start(warmUpTime, recordPeriod int) <-chan Snapshot {
	ch := make(chan Snapshot)

	go t.work(warmUpTime, recordPeriod, ch)

	return ch
}

func (t *App) work(M, N int, ch chan Snapshot) {
	defer close(ch)
	t.logg.Debug("[APP] started... Hello!")

	var cpuCh chan apptop.Cpu
	if t.conf.Metrics.Cpu {
		topRunner := apptop.New(t.conf.App.TopPath, apptop.NewParser(), t.logg)
		cpuCh = topRunner.Run(t.ctx, M, N)
	} else {
		t.logg.Debug("[APP] CPU metric is disabled")
	}

	var disksIOCh chan []appiostat.DiskIO
	var disksStatCh chan []appdf.DiskInfo
	if t.conf.Metrics.Disks {
		iostatRunner := appiostat.New(t.conf.App.IostatPath, t.logg, appiostat.NewParser())
		disksIOCh = iostatRunner.Run(t.ctx, M, N)

		dfRunner := appdf.New(t.conf.App.DfPath, t.logg, appdf.NewParser(t.logg))
		disksStatCh = dfRunner.Run(t.ctx, M, N)
	} else {
		t.logg.Debug("[APP] discs metric is disabled")
	}

	var talkersCh chan apptcpdump.TopTalkers
	if t.conf.Metrics.Network {
		tcpDumpRunner := apptcpdump.New(
			t.conf.App.TimeoutPath,
			t.conf.App.TcpDumpPath,
			t.logg,
			apptcpdump.NewParser(t.logg),
		)
		talkersCh = tcpDumpRunner.Run(t.ctx, M, N)
	} else {
		t.logg.Debug("[APP] network metric is disabled")
	}

	var connsCh chan appnetstat.ConnectData
	if t.conf.Metrics.Connections {
		netStatRunner := appnetstat.New(
			t.conf.App.NetStatPath,
			t.logg,
			appnetstat.NewParser(t.logg),
		)
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
			// проверять закрытие каналов
			var cpu apptop.Cpu
			var cpuOpened bool
			if t.conf.Metrics.Cpu {
				t.logg.Debug("[APP] waiting cpu")
				cpu, cpuOpened = <-cpuCh

				if !cpuOpened {
					t.logg.Debug("[APP] cpu channel is closed")
					break L
				}
			}

			var disksIO []appiostat.DiskIO
			var disksIOOpened bool

			var disksInfo []appdf.DiskInfo
			var disksInfoOpened bool
			if t.conf.Metrics.Disks {
				t.logg.Debug("[APP] waiting disks io")
				disksIO, disksIOOpened = <-disksIOCh

				if !disksIOOpened {
					t.logg.Debug("[APP] disks io is closed")
					break L
				}

				t.logg.Debug("[APP] waiting disks stats")
				disksInfo, disksInfoOpened = <-disksStatCh

				if !disksInfoOpened {
					t.logg.Debug("[APP] disks info is closed")
					break L
				}
			}

			var talkers apptcpdump.TopTalkers
			var talkersChOpened bool
			if t.conf.Metrics.Network {
				t.logg.Debug("[APP] waiting connections")
				talkers, talkersChOpened = <-talkersCh

				if !talkersChOpened {
					t.logg.Debug("[APP] network talkers channel is closed")
					break L
				}
			}

			var conns appnetstat.ConnectData
			var connsChOpened bool
			if t.conf.Metrics.Connections {
				t.logg.Debug("[APP] waiting connections")
				conns, connsChOpened = <-connsCh

				if !connsChOpened {
					t.logg.Debug("[APP] connections channel is closed")
					break L
				}
			}

			ch <- Snapshot{
				Cpu:                  cpu,
				DisksIO:              disksIO,
				DisksInfo:            disksInfo,
				ConnectsInfo:         conns.Infos,
				ConnectsStates:       conns.States,
				TopTalkersByProtocol: talkers.ByProtocol,
				TopTalkersByTraffic:  talkers.ByTraffic,
			}
		}
	}

	t.logg.Info("[APP] finished... Good bye!")
}

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

	go t.prepareDataChannels(warmUpTime, recordPeriod, ch)

	return ch
}

func (t *App) prepareDataChannels(warmUpTime, period int, ch chan Snapshot) {
	defer close(ch)
	t.logg.Debug("[APP] started... Hello!")

	var cpuCh chan apptop.CPU
	if t.conf.Metrics.CPU {
		topRunner := apptop.New(t.conf.App.TopPath, apptop.NewParser(t.logg), t.logg)
		cpuCh = topRunner.Run(t.ctx, warmUpTime, period)
	} else {
		t.logg.Debug("[APP] CPU metric is disabled")
	}

	var disksIOCh chan []appiostat.DiskIO
	var disksStatCh chan []appdf.DiskInfo
	if t.conf.Metrics.Disks {
		iostatRunner := appiostat.New(t.conf.App.IostatPath, t.logg, appiostat.NewParser(t.logg))
		disksIOCh = iostatRunner.Run(t.ctx, warmUpTime, period)

		dfRunner := appdf.New(t.conf.App.DfPath, t.logg, appdf.NewParser(t.logg))
		disksStatCh = dfRunner.Run(t.ctx, warmUpTime, period)
	} else {
		t.logg.Debug("[APP] discs metric is disabled")
	}

	var talkersCh chan apptcpdump.TopTalkers
	if t.conf.Metrics.Network {
		tcpDumpRunner := apptcpdump.New(
			t.conf.App.TimeoutPath,
			t.conf.App.TCPDumpPath,
			t.logg,
			apptcpdump.NewParser(t.logg),
		)
		talkersCh = tcpDumpRunner.Run(t.ctx, warmUpTime, period)
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
		connsCh = netStatRunner.Run(t.ctx, warmUpTime, period)
	} else {
		t.logg.Debug("[APP] connection metrics disabled")
	}

L:
	for {
		select {
		case <-t.ctx.Done():
			break L
		default:
			ch <- t.proxySnaphot(cpuCh, disksIOCh, disksStatCh, talkersCh, connsCh)
		}
	}

	t.logg.Info("[APP] finished... Good bye!")
}

func (t *App) proxySnaphot(
	cpuCh chan apptop.CPU,
	disksIOCh chan []appiostat.DiskIO,
	disksStatCh chan []appdf.DiskInfo,
	talkersCh chan apptcpdump.TopTalkers,
	connsCh chan appnetstat.ConnectData,
) Snapshot {
	var cpu apptop.CPU
	var cpuOpened bool

	if t.conf.Metrics.CPU {
		t.logg.Debug("[APP] waiting cpu")
		cpu, cpuOpened = <-cpuCh

		if !cpuOpened {
			t.logg.Warn("[APP] cpu channel is closed")
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
			t.logg.Warn("[APP] disks io is closed")
		}

		t.logg.Debug("[APP] waiting disks stats")
		disksInfo, disksInfoOpened = <-disksStatCh

		if !disksInfoOpened {
			t.logg.Warn("[APP] disks info is closed")
		}
	}

	var talkers apptcpdump.TopTalkers
	var talkersChOpened bool
	if t.conf.Metrics.Network {
		t.logg.Debug("[APP] waiting connections")
		talkers, talkersChOpened = <-talkersCh

		if !talkersChOpened {
			t.logg.Warn("[APP] network talkers channel is closed")
		}
	}

	var conns appnetstat.ConnectData
	var connsChOpened bool
	if t.conf.Metrics.Connections {
		t.logg.Debug("[APP] waiting connections")
		conns, connsChOpened = <-connsCh

		if !connsChOpened {
			t.logg.Warn("[APP] connections channel is closed")
		}
	}

	return Snapshot{
		CPU:                  cpu,
		DisksIO:              disksIO,
		DisksInfo:            disksInfo,
		ConnectsInfo:         conns.Infos,
		ConnectsStates:       conns.States,
		TopTalkersByProtocol: talkers.ByProtocol,
		TopTalkersByTraffic:  talkers.ByTraffic,
	}
}

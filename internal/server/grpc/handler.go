package internalgrpc

import (
	"context"
	"fmt"
	"log"

	"github.com/Cranky4/go-top/api/TopService"
	pb "github.com/Cranky4/go-top/api/TopService"
	"github.com/Cranky4/go-top/internal/app"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler = pb.TopServiceServer

type handler struct {
	pb.UnimplementedTopServiceServer
	ctx  context.Context
	app  *app.App
	logg *log.Logger
}

func NewHandler(ctx context.Context, app *app.App, logger *log.Logger) (Handler, error) {
	return &handler{app: app, logg: logger}, nil
}

func (h *handler) StreamSnapshots(r *TopService.SnapshotRequest, srv TopService.TopService_StreamSnapshotsServer) error {
	h.logg.Printf("Client connected with params: m=%d, n=%d", r.M, r.N)
	ch := h.app.Start(r.M, r.N)

	for {
		select {
		case <-srv.Context().Done():
			h.logg.Printf("Client disconnected")
			return nil
		case s := <-ch:
			disksIO := make([]*pb.DiskIO, 0, len(s.DisksIO))
			for _, d := range s.DisksIO {
				disksIO = append(disksIO, &pb.DiskIO{
					Device: d.Device,
					Tps:    d.Tps,
					Kbps:   d.Kbps,
				})
			}

			disksInfo := make([]*pb.DiskInfo, 0, len(s.DisksInfo))
			for _, d := range s.DisksInfo {
				disksInfo = append(disksInfo, &pb.DiskInfo{
					Name:            d.Name,
					UsedBytes:       int64(d.UsedBytes),
					AvailableBytes:  int64(d.AvailableBytes),
					UsageBytes:      fmt.Sprintf("%d%%", d.UsageBytes),
					UsedInodes:      int64(d.UsedInodes),
					AvailableInodes: int64(d.AvailableInodes),
					UsageInodes:     fmt.Sprintf("%d%%", d.UsageInodes),
				})
			}

			snapshot := pb.Snapshot{
				StartTime: &timestamppb.Timestamp{
					Seconds: s.StartTime.Unix(),
					Nanos:   int32(s.StartTime.Nanosecond()),
				},
				FinishTime: &timestamppb.Timestamp{
					Seconds: s.FinishTime.Unix(),
					Nanos:   int32(s.FinishTime.Nanosecond()),
				},
				Cpu: &pb.Cpu{
					Avg: &pb.CpuAvg{
						Min:     s.Cpu.Avg.Min,
						Five:    s.Cpu.Avg.Five,
						Fifteen: s.Cpu.Avg.Fifteen,
					},
					State: &pb.CpuState{
						User:   s.Cpu.State.User,
						System: s.Cpu.State.System,
						Idle:   s.Cpu.State.Idle,
					},
				},
				DisksIO:   disksIO,
				DisksInfo: disksInfo,
				// 	DisksInfo: []*pb.DiskInfo{
				// 		{
				// 			Name:            "/dev/nwme0n1",
				// 			UsedBytes:       2414232,
				// 			AvailableBytes:  42423,
				// 			UsageBytes:      "38%",
				// 			UsedInodes:      4123,
				// 			AvailableInodes: 232,
				// 			UsageInodes:     "21%",
				// 		},
				// 		{
				// 			Name:            "/dev/nwme0n2",
				// 			UsedBytes:       2414232,
				// 			AvailableBytes:  42423,
				// 			UsageBytes:      "38%",
				// 			UsedInodes:      4123,
				// 			AvailableInodes: 232,
				// 			UsageInodes:     "21%",
				// 		},
				// 	},
				// 	TopTalkersByProtocol: []*pb.TopTalkerByProtocol{
				// 		{
				// 			Protocol: "UDP",
				// 			Bytes:    24232,
				// 			Percent:  "73%",
				// 		},
				// 		{
				// 			Protocol: "TCP",
				// 			Bytes:    8432,
				// 			Percent:  "25%",
				// 		},
				// 	},
				// 	TopTalkersByTraffic: []*pb.TopTalkerByTraffic{
				// 		{
				// 			Source:         "127.0.0.1",
				// 			Destination:    "127.0.0.2",
				// 			Protocol:       "UDP",
				// 			BytesPerSecond: 174,
				// 		},
				// 		{
				// 			Source:         "127.0.0.1",
				// 			Destination:    "127.0.0.2",
				// 			Protocol:       "TCP",
				// 			BytesPerSecond: 23,
				// 		},
				// 	},
				// 	ConnectsInfo: []*pb.ConnectInfo{
				// 		{
				// 			Command:  "ping",
				// 			Pid:      2312,
				// 			User:     "root",
				// 			Protocol: "UDP",
				// 			Port:     ":90",
				// 		},
				// 		{
				// 			Command:  "smtth",
				// 			Pid:      2344,
				// 			User:     "root",
				// 			Protocol: "TCP",
				// 			Port:     ":9230",
				// 		},
				// 	},
				// 	ConnectsStates: []*pb.ConnectState{
				// 		{
				// 			Protocol: "UDP",
				// 			State:    "READY",
				// 		},
				// 		{
				// 			Protocol: "TCP",
				// 			State:    "BAC",
				// 		},
				// 	},
			}

			if err := srv.Send(&snapshot); err != nil {
				return err
			}
		}
	}
}

package internalgrpc

import (
	"context"
	"log"

	"github.com/Cranky4/go-top/api/TopService"
	pb "github.com/Cranky4/go-top/api/TopService"
	"github.com/Cranky4/go-top/internal/app"
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
	ch := h.app.Start(int(r.M), int(r.N))

	for {
		select {
		case <-srv.Context().Done():
			h.logg.Printf("Client disconnected")
			return nil
		case s, opened := <-ch:
			if !opened {
				return nil
			}

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
					UsedBytes:       int32(d.UsedBytes),
					AvailableBytes:  int32(d.AvailableBytes),
					UsageBytes:      int32(d.UsageBytes),
					UsedInodes:      int32(d.UsedInodes),
					AvailableInodes: int32(d.AvailableInodes),
					UsageInodes:     int32(d.UsageInodes),
				})
			}

			connectsInfo := make([]*pb.ConnectInfo, 0, len(s.ConnectsInfo))
			for _, v := range s.ConnectsInfo {
				connectsInfo = append(connectsInfo, &pb.ConnectInfo{
					Command:  v.Command,
					Pid:      int32(v.Pid),
					Protocol: v.Protocol,
					Port:     int32(v.Port),
				})
			}

			connectsStates := make([]*pb.ConnectState, 0, len(s.ConnectsStates))
			for _, v := range s.ConnectsStates {
				connectsStates = append(connectsStates, &pb.ConnectState{
					Protocol: v.Protocol,
					State:    v.State,
				})
			}

			topTalkersByProtocol := make([]*pb.TopTalkerByProtocol, 0, len(s.TopTalkersByProtocol))
			for _, v := range s.TopTalkersByProtocol {
				topTalkersByProtocol = append(topTalkersByProtocol, &pb.TopTalkerByProtocol{
					Protocol: v.Protocol,
					Bytes:    int32(v.Bytes),
					Percent:  v.Percent,
				})
			}

			topTalkersByTraffic := make([]*pb.TopTalkerByTraffic, 0, len(s.TopTalkersByTraffic))
			for _, v := range s.TopTalkersByTraffic {
				topTalkersByTraffic = append(topTalkersByTraffic, &pb.TopTalkerByTraffic{
					Source:         v.Source,
					Destination:    v.Destination,
					Protocol:       v.Protocol,
					BytesPerSecond: v.BytesPerSecond,
				})
			}

			snapshot := pb.Snapshot{
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
				DisksIO:              disksIO,
				DisksInfo:            disksInfo,
				ConnectsInfo:         connectsInfo,
				ConnectsStates:       connectsStates,
				TopTalkersByProtocol: topTalkersByProtocol,
				TopTalkersByTraffic:  topTalkersByTraffic,
			}

			if err := srv.Send(&snapshot); err != nil {
				return err
			}
		}
	}
}

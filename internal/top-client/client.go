package topclient

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	pb "github.com/Cranky4/go-top/api/TopService"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TopClient struct {
	ctx  context.Context
	conf Config
	logg Logger
}

func New(ctx context.Context, conf Config, logg Logger) *TopClient {
	logg.Debug(fmt.Sprintf("[Config] %#v", conf))
	return &TopClient{
		ctx:  ctx,
		conf: conf,
		logg: logg,
	}
}

func (c *TopClient) Start() error {
	// TODO: creds?
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())

	conn, err := grpc.DialContext(c.ctx, c.conf.Grpc.Addr, opts)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewTopServiceClient(conn)

	stream, err := client.StreamSnapshots(c.ctx, &pb.SnapshotRequest{
		M: uint32(c.conf.Client.M),
		N: uint32(c.conf.Client.N),
	})
	if err != nil {
		return err
	}

L:
	for {
		select {
		case <-c.ctx.Done():
			break L
		default:
			snapshot, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break L
			}
			if err != nil {
				return err
			}

			c.ProcessSnapshot(snapshot)
		}
	}

	c.logg.Info("top client stopped")

	return nil
}

func (c *TopClient) ProcessSnapshot(sn *pb.Snapshot) {
	var disksIO strings.Builder
	defer disksIO.Reset()
	for _, d := range sn.DisksIO {
		disksIO.WriteString(fmt.Sprintf("%s - %.02f kbps, %.02f tps\n", d.Device, d.Kbps, d.Tps))
	}

	var disksInfo strings.Builder
	defer disksInfo.Reset()
	for _, d := range sn.DisksInfo {
		disksInfo.WriteString(
			fmt.Sprintf("%s - %d bytes used, %d bytes available - %d%%\n", d.Name, d.UsedBytes, d.AvailableBytes, d.UsageBytes),
		)
	}

	var topTalkersByProtocol strings.Builder
	defer topTalkersByProtocol.Reset()
	for _, t := range sn.TopTalkersByProtocol {
		topTalkersByProtocol.WriteString(
			fmt.Sprintf("%s %d %s\n", t.Protocol, t.Bytes, t.Percent),
		)
	}

	var topTalkersByTraffic strings.Builder
	defer topTalkersByTraffic.Reset()
	for _, t := range sn.TopTalkersByTraffic {
		topTalkersByTraffic.WriteString(
			fmt.Sprintf("%s %s >%s %s %.02f\n", t.Protocol, t.Source, t.Destination, t.Protocol, t.BytesPerSecond),
		)
	}

	var connectsInfo strings.Builder
	defer connectsInfo.Reset()
	for _, t := range sn.ConnectsInfo {
		connectsInfo.WriteString(
			fmt.Sprintf("%s [%d] %s %d\n", t.Command, t.Pid, t.Protocol, t.Port),
		)
	}

	var connectsStates strings.Builder
	defer connectsStates.Reset()
	for _, t := range sn.ConnectsStates {
		connectsStates.WriteString(
			fmt.Sprintf("%s - %s\n", t.Protocol, t.State),
		)
	}

	c.logg.Info(
		fmt.Sprintf(
			"\nLoad average: %.02f %.02f %.02f\nStates: %.02fus %.02fsy %.02fid\n\nDisks IO:\n%s\n\nDisks Info:\n%s\n\n"+
				"Top Talkers by protocol:\n%s\n\nTop Talkers by traffic:\n%s\n\nConnections:\n%s\n\nProtocol states:\n%s\n",
			sn.Cpu.Avg.Min, sn.Cpu.Avg.Five, sn.Cpu.Avg.Fifteen,
			sn.Cpu.State.User, sn.Cpu.State.System, sn.Cpu.State.Idle,
			disksIO.String(),
			disksInfo.String(),
			topTalkersByProtocol.String(),
			topTalkersByTraffic.String(),
			connectsInfo.String(),
			connectsStates.String(),
		),
	)
}

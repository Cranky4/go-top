package internalgrpc

import (
	"log"
	"time"

	"github.com/Cranky4/go-top/api/TopService"
	pb "github.com/Cranky4/go-top/api/TopService"
	"github.com/Cranky4/go-top/internal/top"
)

type Handler = pb.TopServiceServer

type handler struct {
	pb.UnimplementedTopServiceServer
	app  *top.Top
	logg *log.Logger
}

func NewHandler(app *top.Top, logger *log.Logger) (Handler, error) {
	return &handler{app: app, logg: logger}, nil
}

func (h *handler) StreamSnapshots(r *TopService.SnapshotRequest, srv TopService.TopService_StreamSnapshotsServer) error {
	h.logg.Printf("%#v", r)

	for {
		time.Sleep(10 * time.Second)
		snapshot := pb.Snapshot{}

		srv.Send(&snapshot)
	}
}

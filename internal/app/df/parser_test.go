package appdf

import (
	"testing"

	"github.com/Cranky4/go-top/internal/logger"
	"github.com/stretchr/testify/require"
)

func TestParseInodes(t *testing.T) {
	logg := logger.New("error", 0)
	parser := NewParser(logg)

	t.Run("succeeded parsed df inodes output", func(t *testing.T) {
		in := "Filesystem        Inodes   IUsed    IFree IUse% Mounted on\n" +
			"/dev/nvme0n1p8         0       0        0     - /etc/hosts\n" +
			"/dev/nvme0n1p10 18612224 1451647 17160577    8% /opt/app/logs\n"

		ex := []DiskInfo{
			{
				Name:            "/dev/nvme0n1p8",
				AvailableInodes: 0,
				UsedInodes:      0,
				UsageInodes:     0,
			},
			{
				Name:            "/dev/nvme0n1p10",
				AvailableInodes: 18612224,
				UsedInodes:      1451647,
				UsageInodes:     8,
			},
		}

		out := parser.ParseInodes(in)
		require.Equal(t, ex, out)
	})

	t.Run("error input", func(t *testing.T) {
		in := "Filesystem      1K-blocks      Used Available Use% Mounted on\n" +
			"/dev/nvme0n1p8   65339392  48946056  11318152  z% /etc/hosts\n"

		out := parser.ParseBytes(in)
		require.Empty(t, out)
	})

	t.Run("succeeded parsed df bytes output", func(t *testing.T) {
		in := "Filesystem      1K-blocks      Used Available Use% Mounted on\n" +
			"/dev/nvme0n1p8   65339392  48946056  11318152  82% /etc/hosts\n" +
			"/dev/nvme0n1p10 291902464 112315688 164685892  41% /opt/app/logs\n"

		ex := []DiskInfo{
			{
				Name:           "/dev/nvme0n1p8",
				AvailableBytes: 11318152,
				UsedBytes:      48946056,
				UsageBytes:     82,
			},
			{
				Name:           "/dev/nvme0n1p10",
				AvailableBytes: 164685892,
				UsedBytes:      112315688,
				UsageBytes:     41,
			},
		}

		out := parser.ParseBytes(in)
		require.Equal(t, ex, out)
	})

	t.Run("error input", func(t *testing.T) {
		in := "Filesystem      1K-blocks      Used Available Use% Mounted on\n" +
			"/dev/nvme0n1p8   65339392  48946056  11318152  8z2% /etc/hosts\n"

		out := parser.ParseBytes(in)
		require.Empty(t, out)
	})
}

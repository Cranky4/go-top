package appdf

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseInodes(t *testing.T) {
	parser := NewParser(nil)

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

		out, err := parser.ParseInodes(in)

		require.Nil(t, err)
		require.Equal(t, ex, out)
	})

	t.Run("error input", func(t *testing.T) {
		in := "Filesystem      1K-blocks      Used Available Use% Mounted on\n" +
			"/dev/nvme0n1p8   65339392  48946056  11318152  z% /etc/hosts\n" +
			"/dev/nvme0n1p10 291902464 112315688 164685892  41% /opt/app/logs\n"

		out, err := parser.ParseBytes(in)

		require.NotNil(t, err)
		var expErr *ErrCannotParseInput
		require.True(t, errors.As(err, &expErr))
		require.Nil(t, out)
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

		out, err := parser.ParseBytes(in)

		require.Nil(t, err)
		require.Equal(t, ex, out)
	})

	t.Run("error input", func(t *testing.T) {
		in := "Filesystem      1K-blocks      Used Available Use% Mounted on\n" +
			"/dev/nvme0n1p8   65339392  48946056  11318152  8z2% /etc/hosts\n" +
			"/dev/nvme0n1p10 291902464 112315688 164685892  41% /opt/app/logs\n"

		out, err := parser.ParseBytes(in)

		require.NotNil(t, err)
		var expErr *ErrCannotParseInput
		require.True(t, errors.As(err, &expErr))
		require.Nil(t, out)
	})
}

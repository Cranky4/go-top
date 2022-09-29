package apptcpdump

import (
	"testing"
	"time"

	"github.com/Cranky4/go-top/internal/logger"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	logg := logger.New("error", 0)

	parser := NewParser(logg)

	t.Run("succeeded parsed", func(t *testing.T) {
		in := "tcpdump: verbose output suppressed, use -v or -vv for full protocol decode\n" +
			"listening on any, link-type LINUX_SLL (Linux cooked), capture size 262144 bytes\n" +
			"2022-09-29 17:27:44.357727 IP6 fe80::42:c1ff:fe1d:f67d.5353 > ff02::fb.5353: UDP, length 118\n" +
			"2022-09-29 17:27:45.023316 IP6 fe80::1c1b:f3ff:fe4d:b90c.5353 > ff02::fb.5353: UDP, length 118\n" +
			"\n" +
			"2 packets captured\n" +
			"2 packets received by filter\n" +
			"0 packets dropped by kernel\n"

		time1, err := time.Parse("2006-01-02 15:04:05.999999999", "2022-09-29 17:27:44.357727")
		require.Nil(t, err)

		time2, err := time.Parse("2006-01-02 15:04:05.999999999", "2022-09-29 17:27:45.023316")
		require.Nil(t, err)

		ex := []TcpDumpLine{
			{
				Time:        time1,
				Type:        "IP6",
				Source:      "fe80::42:c1ff:fe1d:f67d.5353",
				Destination: "ff02::fb.5353",
				Protocol:    "UDP",
				Bytes:       118,
			},

			{
				Time:        time2,
				Type:        "IP6",
				Source:      "fe80::1c1b:f3ff:fe4d:b90c.5353",
				Destination: "ff02::fb.5353",
				Protocol:    "UDP",
				Bytes:       118,
			},
		}

		out, err := parser.Parse(in)

		require.Nil(t, err)
		require.Equal(t, ex, out)
	})

	t.Run("empty input", func(t *testing.T) {
		in := "tcpdump: verbose output suppressed, use -v or -vv for full protocol decode\n" +
			"listening on any, link-type LINUX_SLL (Linux cooked), capture size 262144 bytes\n" +
			"\n" +
			"0 packets captured\n" +
			"0 packets received by filter\n" +
			"0 packets dropped by kernel\n"

		out, err := parser.Parse(in)

		require.Nil(t, err)
		require.Empty(t, out)
	})

	t.Run("invalid input", func(t *testing.T) {
		in := "Linux 5.18.17-1-MANJARO (7cc95de5443a)  09/27/22        _x86_64_        (16 CPU)\n" +
			"\n" +
			"Device             tps    kB_read/s    kB_wrtn/s    kB_read    kB_wrtn\n" +
			"nvme0n1          72.84    33    329.34      2369.18    594383.7  244  42758633\n"

		out, err := parser.Parse(in)

		require.Nil(t, err)
		require.Empty(t, out)
	})
}

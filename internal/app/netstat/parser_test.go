package appnetstat

import (
	"testing"

	"github.com/Cranky4/go-top/internal/logger"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	logg := logger.New("error", 0)
	parser := NewParser(logg)

	t.Run("succeeded parsed", func(t *testing.T) {
		in := "Active Internet connections (servers and established)\n" +
			"Proto Recv-Q Send-Q Local Address           Foreign Address         State       PID/Program name    \n" +
			"tcp        0      0 127.0.0.11:36269        0.0.0.0:*               LISTEN      -                   \n" +
			"tcp6       0      0 :::9990                 :::*                    LISTEN      7/top    \n"

		ex := []NetStatRow{
			{
				Proto:       "tcp",
				RecvQ:       0,
				SendQ:       0,
				LocalAddr:   "127.0.0.11:36269",
				LocalPort:   36269,
				ForeignAddr: "0.0.0.0:*",
				ForeignPort: 0,
				State:       "LISTEN",
				PID:         0,
				Programm:    "",
			},
			{
				Proto:       "tcp6",
				RecvQ:       0,
				SendQ:       0,
				LocalAddr:   ":::9990",
				LocalPort:   9990,
				ForeignAddr: ":::*",
				ForeignPort: 0,
				State:       "LISTEN",
				PID:         7,
				Programm:    "top",
			},
		}

		out := parser.Parse(in)
		require.Equal(t, ex, out)
	})

	t.Run("invalid input", func(t *testing.T) {
		in := "Linux 5.18.17-1-MANJARO (7cc95de5443a)  09/27/22        _x86_64_        (16 CPU)\n" +
			"\n" +
			"Device             tps    kB_read/s    kB_wrtn/s    kB_read    kB_wrtn\n" +
			"nvme0n1          72.84    33    329.34      2369.18    594383.7  244  42758633\n"

		out := parser.Parse(in)
		require.Empty(t, out)
	})
}

package apptop

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	parser := NewParser()

	t.Run("succeeded parsed", func(t *testing.T) {
		in := "top - 20:06:17 up  5:39,  0 users,  load average: 0.60, 0.56, 0.62\n" +
			"Tasks:   4 total,   1 running,   3 sleeping,   0 stopped,   0 zombie\n" +
			"%Cpu(s):  1.9 us,  1.1 sy,  0.4 ni, 96.4 id,  0.1 wa,  0.0 hi,  0.2 si,  0.0 st\n" +
			"KiB Mem : 15715104 total,  2289212 free,  5072812 used,  8353080 buff/cache\n" +
			"KiB Swap: 16777212 total, 16777212 free,        0 used. 10152312 avail Mem\n" +
			"\n" +
			"PID USER      PR  NI    VIRT    RES    SHR S  %CPU %MEM     TIME+ COMMAND\n" +
			"1 root      20   0    4640    856    792 S   0.0  0.0   0:00.01 sh\n" +
			"7 root      20   0  722724  21488   8296 S   0.0  0.1   0:18.81 top\n" +
			"16773 root      20   0   18520   3536   3120 S   0.0  0.0   0:00.01 bash\n" +
			"18934 root      20   0   36500   3144   2788 R   0.0  0.0   0:00.00 top\n"

		ex := Cpu{
			Avg: CpuAvg{
				Min:     0.6,
				Five:    0.56,
				Fifteen: 0.62,
			},
			State: CpuState{
				User:   1.9,
				System: 1.1,
				Idle:   96.4,
			},
		}

		out, err := parser.Parse(in)

		require.Nil(t, err)
		require.Equal(t, ex, out)
	})

	t.Run("invalid input", func(t *testing.T) {
		in := "Linux 5.18.17-1-MANJARO (7cc95de5443a)  09/27/22        _x86_64_        (16 CPU)\n" +
			"\n" +
			"Device             tps    kB_read/s    kB_wrtn/s    kB_read    kB_wrtn\n" +
			"nvme0n1          72.84    33    329.34      2369.18    594383.7  244  42758633\n"

		out, err := parser.Parse(in)

		require.NotNil(t, err)
		var expErr *ErrCannotParseInput
		require.True(t, errors.As(err, &expErr))
		require.Empty(t, out)
	})
}

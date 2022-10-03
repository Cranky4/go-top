package appiostat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	parser := NewParser(&TLogger{})

	t.Run("succeeded parsed", func(t *testing.T) {
		in := "Linux 5.18.17-1-MANJARO (7cc95de5443a)  09/27/22        _x86_64_        (16 CPU)\n" +
			"\n" +
			"Device             tps    kB_read/s    kB_wrtn/s    kB_read    kB_wrtn\n" +
			"nvme0n1          72.84       329.34      2369.18    5943837   42758633\n"

		ex := []IostatRow{
			{
				Device:    "nvme0n1",
				Tps:       72.84,
				KbpsRead:  329.34,
				KbpsWrite: 2369.18,
				KbRead:    5943837,
				KbWrite:   42758633,
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
		require.Nil(t, out)
	})
}

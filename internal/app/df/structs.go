package appdf

import "fmt"

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type DiskInfo struct {
	Name            string // /dev/nvme0n1p3     df -k
	UsedBytes       int    // 39131452
	AvailableBytes  int    // 64169664
	UsageBytes      int    // 38%
	UsedInodes      int    // 272911         df -i
	AvailableInodes int    // 64224433
	UsageInodes     int    // 1%
}

type Parser interface {
	ParseBytes(in string) []DiskInfo
	ParseInodes(in string) []DiskInfo
}

type ErrCannotParseInput struct {
	Input string
}

func (e *ErrCannotParseInput) Error() string {
	return fmt.Sprintf("cannot parse %s", e.Input)
}

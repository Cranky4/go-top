package appiostat

type DiskIO struct {
	Device    string  // nvme0n1
	Tps, Kbps float32 // 52.86, 665.63 + 780.25
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

package topclient

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Config struct {
	Client ClientConf
	Logg   LoggerConf
	Grpc   GrpcConf
}

type ClientConf struct {
	WarmingUpTime, SnapshotPeriod int
}

type LoggerConf struct {
	Level string
}

type GrpcConf struct {
	Addr string
}

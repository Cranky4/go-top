package top

type Config struct {
	Top  TopConf
	Logg LoggerConf
	Grpc GrpcConf
}

type TopConf struct {
	Metrics string
}

type LoggerConf struct {
	Level string
}

type GrpcConf struct {
	Addr, RequestLogFile string
}

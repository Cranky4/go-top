package top

type Top struct {
	conf Config
	logg Logger
}

func NewTop(conf Config, logg Logger) *Top {
	return &Top{conf: conf, logg: logg}
}

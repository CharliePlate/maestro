package maestro

type QueueConfig struct{}

type Queue struct {
	Container Container
	Writer    ContainerWriter
	Cfg       QueueConfig
	Name      string
}

type ContainerWriter interface {
	Write(item QueueItem) error
}

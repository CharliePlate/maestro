package maestro

import "io"

type QueueConfig struct{}

type Queue struct {
	Container Container
	Watcher   Watcher
	Cfg       QueueConfig
	Name      string
}

type Acknowledger interface {
	Acknowledge(id string) error
}

type Subscriber interface {
	Subscribe(queue string) error
}

type Unsubscriber interface {
	Unsubscribe(queue string) error
}

// The writer and reader interface will be what is written/read from the net connection.
type Peer interface {
	io.Writer
	io.Reader
	Acknowledger
	Subscriber
	Unsubscriber
}

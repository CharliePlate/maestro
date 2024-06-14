package maestro

import (
	"log/slog"
)

type Config struct {
	Logger slog.Logger
}

type Maestro struct {
	Config          Config
	QueueRegisterer QueueRegisterer
	Server          Server
	Queues          []Queue
}

type QueueRegisterer interface {
	Register(q Queue) error
}

func (m *Maestro) WithQueue(q Queue) error {
	if err := m.QueueRegisterer.Register(q); err != nil {
		return err
	}

	return nil
}

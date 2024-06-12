package maestro

import (
	"log/slog"
)

type Config struct {
	Logger slog.Logger
}

type Maestro struct {
	Config Config
	Queues []Queue
}

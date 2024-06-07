package maestro

import "context"

type Watcher interface {
	Watch(ctx context.Context, c chan QueueUpdateMessage) error
}

package maestro

import "context"

// Test Go Change

type Watcher interface {
	Watch(ctx context.Context, c chan QueueUpdateMessage) error
}

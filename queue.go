package maestro

import "context"

type Queue struct {
	Container Container
	Handler   QueueHandler
	Receiver  Receiver
	Name      string
}

func (q *Queue) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch := make(chan *QueueUpdateMessage)
	errCh := make(chan error)
	defer close(ch)
	defer close(errCh)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				q.Receiver.Listen(ctx, ch, errCh)
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-ch:
			q.Handler.Handle(msg, errCh)
		case err := <-errCh:
			// handle error
			_ = err
		}
	}
}

type QueueUpdateMessage struct {
	Content any
	MsgID   string
	OpType  OpType
}

func (qu *QueueUpdateMessage) ID() string {
	return qu.MsgID
}

func (qu *QueueUpdateMessage) Data() interface{} {
	return qu.Content
}

func (qu *QueueUpdateMessage) SetID(id string) {
	qu.MsgID = id
}

func (qu *QueueUpdateMessage) SetData(data interface{}) {
	qu.Content = data
}

type QueueHandler interface {
	Handle(m *QueueUpdateMessage, errChan chan error)
	SetContainer(c Container)
	Container() *Container
}

type Receiver interface {
	Listen(ctx context.Context, ch chan *QueueUpdateMessage, errCh chan error)
}

type QueueItem interface {
	ID() string
	Data() interface{}
	SetID(id string)
	SetData(data interface{})
}

type Container interface {
	Push(item QueueItem)
	Pop() (QueueItem, error)
	Len() int
	Items() []QueueItem
	Find(id string) (QueueItem, error)
	Delete(id string) error
}

type OpType int

const (
	OpTypeInsert OpType = iota
	OpTypeUpdate
	OpTypeDelete
)

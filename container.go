package maestro

import "errors"

type OpType int

const (
	OpTypeInsert OpType = iota
	OpTypeUpdate
	OpTypeDelete
)

type QueueUpdateMessage struct {
	Data   interface{}
	ID     string
	OpType OpType
}

type QueueItem interface {
	ID() string
	Data() interface{}
}

type Container interface {
	Push(item QueueItem)
	Pop() (QueueItem, error)
	Len() int
	Items() ([]QueueItem, error)
	Find(id string) (QueueItem, error)
}

var (
	ErrQueueEmpty   = errors.New("queue is empty")
	ErrItemNotFound = errors.New("item not found")
)

type SliceContainer struct {
	Elements []QueueItem
}

func NewSliceContainer() *SliceContainer {
	return &SliceContainer{
		Elements: make([]QueueItem, 0),
	}
}

func (sc *SliceContainer) Push(item QueueItem) {
	sc.Elements = append(sc.Elements, item)
}

func (sc *SliceContainer) Pop() (QueueItem, error) {
	if len(sc.Elements) == 0 {
		return nil, ErrQueueEmpty
	}

	e := sc.Elements[0]

	if len(sc.Elements) == 1 {
		sc.Elements = make([]QueueItem, 0)
	} else {
		sc.Elements = sc.Elements[1:]
	}

	return e, nil
}

func (sc *SliceContainer) Len() int {
	return len(sc.Elements)
}

func (sc *SliceContainer) Items() ([]QueueItem, error) {
	if len(sc.Elements) == 0 {
		return []QueueItem{}, ErrQueueEmpty
	}

	return sc.Elements, nil
}

func (sc *SliceContainer) Find(id string) (QueueItem, error) {
	for _, item := range sc.Elements {
		if item.ID() == id {
			return item, nil
		}
	}
	return nil, ErrItemNotFound
}

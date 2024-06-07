package maestro

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
	Pop() QueueItem
	Len() int
	Items() []QueueItem
	Find(id string) QueueItem
}

type SliceContainer struct {
	items []QueueItem
}

func NewSliceContainer() *SliceContainer {
	return &SliceContainer{
		items: make([]QueueItem, 0),
	}
}

func (sc *SliceContainer) Push(item QueueItem) {
	sc.items = append(sc.items, item)
}

func (sc *SliceContainer) Pop() QueueItem {
	e := sc.items[0]
	sc.items = sc.items[1:]
	return e
}

func (sc *SliceContainer) Len() int {
	return len(sc.items)
}

func (sc *SliceContainer) Items() []QueueItem {
	return sc.items
}

func (sc *SliceContainer) Find(id string) QueueItem {
	for _, item := range sc.items {
		if item.ID() == id {
			return item
		}
	}
	return nil
}

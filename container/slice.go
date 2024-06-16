package container

import (
	"errors"
	"fmt"

	"github.com/charlieplate/maestro"
)

var (
	ErrQueueEmpty   = errors.New("queue is empty")
	ErrItemNotFound = errors.New("item not found")
)

type SliceContainer struct {
	Elements []maestro.QueueItem
}

func NewSliceContainer() *SliceContainer {
	return &SliceContainer{
		Elements: make([]maestro.QueueItem, 0),
	}
}

func (sc *SliceContainer) Push(item maestro.QueueItem) {
	sc.Elements = append(sc.Elements, item)
}

func (sc *SliceContainer) Pop() (maestro.QueueItem, error) {
	if len(sc.Elements) == 0 {
		return nil, ErrQueueEmpty
	}

	e := sc.Elements[0]

	if len(sc.Elements) == 1 {
		sc.Elements = make([]maestro.QueueItem, 0)
	} else {
		sc.Elements = sc.Elements[1:]
	}

	return e, nil
}

func (sc *SliceContainer) Len() int {
	return len(sc.Elements)
}

func (sc *SliceContainer) Items() []maestro.QueueItem {
	// this should never happen but... just in case
	if sc.Elements == nil {
		return []maestro.QueueItem{}
	}

	return sc.Elements
}

func (sc *SliceContainer) Find(id string) (maestro.QueueItem, error) {
	for _, item := range sc.Elements {
		if item.ID() == id {
			return item, nil
		}
	}
	return nil, ErrItemNotFound
}

func (sc *SliceContainer) Delete(id string) error {
	idx, err := sc.findIndexByID(id)
	if err != nil {
		return fmt.Errorf("error deleting id %s: %w", id, err)
	}

	sc.Elements = append(sc.Elements[:idx], sc.Elements[idx+1:]...)
	return nil
}

func (sc *SliceContainer) findIndexByID(id string) (int, error) {
	for i, item := range sc.Elements {
		if item.ID() == id {
			return i, nil
		}
	}

	return -1, ErrItemNotFound
}

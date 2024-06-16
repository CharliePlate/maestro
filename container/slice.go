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
	elements []maestro.QueueItem
	start    int
}

func NewSliceContainer() *SliceContainer {
	return &SliceContainer{
		elements: make([]maestro.QueueItem, 0),
		start:    0,
	}
}

func (sc *SliceContainer) Push(item maestro.QueueItem) {
	sc.elements = append(sc.elements, item)
}

func (sc *SliceContainer) Pop() (maestro.QueueItem, error) {
	if sc.Len() == 0 {
		return nil, ErrQueueEmpty
	}

	e := sc.elements[sc.start]
	sc.start++

	if sc.start > len(sc.elements)/2 {
		sc.elements = sc.elements[sc.start:]
		sc.start = 0
	}

	return e, nil
}

func (sc *SliceContainer) Len() int {
	return len(sc.elements) - sc.start
}

func (sc *SliceContainer) Items() []maestro.QueueItem {
	if sc.elements == nil {
		return []maestro.QueueItem{}
	}
	return sc.elements[sc.start:]
}

func (sc *SliceContainer) Find(id string) (maestro.QueueItem, error) {
	for _, item := range sc.elements {
		if item.ID() == id {
			return item, nil
		}
	}
	return nil, fmt.Errorf("could not find element: %s: %w", id, ErrItemNotFound)
}

func (sc *SliceContainer) Delete(id string) error {
	idx, err := sc.findIndexByID(id)
	if err != nil {
		return fmt.Errorf("error deleting id %s: %w", id, err)
	}

	idx += sc.start

	sc.elements = append(sc.elements[:idx], sc.elements[idx+1:]...)
	return nil
}

func (sc *SliceContainer) findIndexByID(id string) (int, error) {
	for i, item := range sc.elements {
		if item.ID() == id {
			return i, nil
		}
	}

	return -1, ErrItemNotFound
}

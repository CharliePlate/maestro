package maestro_test

import (
	"fmt"
	"testing"

	"github.com/charlieplate/maestro"
	"github.com/stretchr/testify/require"
)

func TestNewSliceContainer(t *testing.T) {
	tests := []struct {
		want *maestro.SliceContainer
		name string
	}{
		{
			name: "Implements Container Interface",
			want: &maestro.SliceContainer{
				Elements: []maestro.QueueItem{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, maestro.NewSliceContainer())
			require.Implements(t, (*maestro.Container)(nil), maestro.NewSliceContainer())
		})
	}
}

type TestQueueItem struct {
	SetID   string
	SetData string
}

func (t *TestQueueItem) ID() string {
	return t.SetID
}

func (t *TestQueueItem) Data() any {
	return t.SetData
}

func TestSliceContainer_Push(t *testing.T) {
	type fields struct {
		elements []maestro.QueueItem
	}
	type args struct {
		item maestro.QueueItem
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []maestro.QueueItem
	}{
		{
			name: "Push Item",
			fields: fields{
				elements: []maestro.QueueItem{},
			},
			args: args{
				item: testQueueItem(0),
			},
			want: []maestro.QueueItem{
				testQueueItem(0),
			},
		},
		{
			name: "Push Item to Non-Empty Container",
			fields: fields{
				elements: []maestro.QueueItem{
					testQueueItem(0),
				},
			},
			args: args{
				item: testQueueItem(1),
			},
			want: makeTestQueueItems(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &maestro.SliceContainer{
				Elements: tt.fields.elements,
			}
			sc.Push(tt.args.item)

			require.Equal(t, tt.want, sc.Elements, "Push() did not add item to container")
		})
	}
}

func TestSliceContainer_Pop(t *testing.T) {
	type fields struct {
		elements []maestro.QueueItem
	}
	tests := []struct {
		want          maestro.QueueItem
		expectedError error
		name          string
		fields        fields
		expectedItems []maestro.QueueItem
	}{
		{
			name: "Pop From Container with 1 Element",
			fields: fields{
				elements: []maestro.QueueItem{
					testQueueItem(0),
				},
			},
			want:          testQueueItem(0),
			expectedItems: []maestro.QueueItem{},
			expectedError: nil,
		},
		{
			name: "Pop From Container with Multiple Elements",
			fields: fields{
				elements: makeTestQueueItems(2),
			},
			want: testQueueItem(0),
			expectedItems: []maestro.QueueItem{
				testQueueItem(1),
			},
			expectedError: nil,
		},
		{
			name: "Pop From Empty Container",
			fields: fields{
				elements: []maestro.QueueItem{},
			},
			want:          nil,
			expectedItems: []maestro.QueueItem{},
			expectedError: maestro.ErrQueueEmpty,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &maestro.SliceContainer{
				Elements: tt.fields.elements,
			}
			item, err := sc.Pop()
			require.Equal(t, tt.expectedItems, sc.Elements, "Unexpected items in container after Pop()")
			require.Equal(t, tt.expectedError, err, "Pop() did not return the expected error")
			require.Equal(t, tt.want, item, "Pop() did not return the expected item")
		})
	}
}

func TestSliceContainer_Len(t *testing.T) {
	type fields struct {
		elements []maestro.QueueItem
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Empty Container",
			fields: fields{
				elements: []maestro.QueueItem{},
			},
			want: 0,
		},
		{
			name: "Container with 1 Element",
			fields: fields{
				elements: []maestro.QueueItem{
					&TestQueueItem{SetID: "testID", SetData: "testData"},
				},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &maestro.SliceContainer{
				Elements: tt.fields.elements,
			}
			if got := sc.Len(); got != tt.want {
				t.Errorf("SliceContainer.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceContainer_Items(t *testing.T) {
	type fields struct {
		elements []maestro.QueueItem
	}
	tests := []struct {
		name   string
		fields fields
		want   []maestro.QueueItem
	}{
		{
			name: "Empty Container",
			fields: fields{
				elements: []maestro.QueueItem{},
			},
			want: []maestro.QueueItem{},
		},
		{
			name: "Container with Data",
			fields: fields{
				elements: makeTestQueueItems(5),
			},
			want: makeTestQueueItems(5),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &maestro.SliceContainer{}

			for _, item := range tt.fields.elements {
				sc.Push(item)
			}

			require.Equal(t, tt.want, sc.Items(), "Items() did not return the expected items")
		})
	}
}

func TestSliceContainer_Find(t *testing.T) {
	type fields struct {
		elements []maestro.QueueItem
	}
	type args struct {
		id string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          maestro.QueueItem
		expectedError error
	}{
		{
			name: "Find Item",
			fields: fields{
				elements: makeTestQueueItems(5),
			},
			args: args{
				id: testQueueItem(3).ID(),
			},
			want:          testQueueItem(3),
			expectedError: nil,
		},
		{
			name: "Find Item Not Found",
			fields: fields{
				elements: makeTestQueueItems(5),
			},
			args: args{
				id: "notFound",
			},
			want:          nil,
			expectedError: maestro.ErrItemNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &maestro.SliceContainer{
				Elements: tt.fields.elements,
			}
			item, err := sc.Find(tt.args.id)
			require.Equal(t, tt.want, item, "Find() did not return the expected item")
			require.Equal(t, tt.expectedError, err, "Find() did not return the expected error")
		})
	}
}

func makeTestQueueItems(count int) []maestro.QueueItem {
	items := []maestro.QueueItem{}

	for i := range count {
		items = append(items,
			&TestQueueItem{SetID: fmt.Sprintf("testId%d", i), SetData: fmt.Sprintf("testData%d", i)},
		)
	}

	return items
}

func testQueueItem(idx int) maestro.QueueItem {
	return &TestQueueItem{SetID: fmt.Sprintf("testId%d", idx), SetData: fmt.Sprintf("testData%d", idx)}
}

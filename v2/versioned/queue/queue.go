package queue

import (
	"sort"
	"sync"

	"github.com/nsu-smp-version-memory/versioned-memory/v2/internal/core"
	"github.com/nsu-smp-version-memory/versioned-memory/v2/internal/timeline"
)

type Queue struct {
	mutex           sync.Mutex
	timeline        *timeline.Timeline[Diff]
	pendingBranches []pendingBranch
	wg              sync.WaitGroup
	merger          timeline.Merger[Diff]
}

func NewQueue() *Queue {
	return &Queue{
		timeline: timeline.NewTimeline[Diff](core.NewSource()),
		merger:   &AppendOrderMerger{},
	}
}

func (q *Queue) Enqueue(value int) {
	q.mutex.Lock()
	q.timeline = q.timeline.NextChange(Diff{Kind: Enqueue, Value: value})
	q.mutex.Unlock()
}

func (q *Queue) Dequeue() {
	q.mutex.Lock()
	q.timeline = q.timeline.NextChange(Diff{Kind: Dequeue})
	q.mutex.Unlock()
}

func (q *Queue) Front() (int, bool) {
	q.mutex.Lock()
	tl := q.timeline
	q.mutex.Unlock()

	data := replayToSlice(tl)
	if len(data) == 0 {
		return 0, false
	}
	return data[0], true
}

func (q *Queue) Items() []int {
	q.mutex.Lock()
	tl := q.timeline
	q.mutex.Unlock()

	data := replayToSlice(tl)

	out := make([]int, len(data))
	copy(out, data)
	return out
}

func (q *Queue) Size() int {
	q.mutex.Lock()
	tl := q.timeline
	q.mutex.Unlock()

	return len(replayToSlice(tl))
}

func (q *Queue) SetMerger(merger timeline.Merger[Diff]) {
	q.mutex.Lock()
	q.merger = merger
	q.mutex.Unlock()
}

func Merge(a, b *Queue) *Queue {
	merger := a.merger

	a.mutex.Lock()
	opsA := a.timeline.Operations()
	a.mutex.Unlock()

	b.mutex.Lock()
	opsB := b.timeline.Operations()
	b.mutex.Unlock()

	result := merger.Merge([][]timeline.Operation[Diff]{opsA, opsB})
	sortOperationsByID(result)

	return &Queue{
		timeline: timeline.FromOperations(core.NewSource(), result),
		merger:   merger,
	}
}

func sortOperationsByID[DIFF any](ops []timeline.Operation[DIFF]) {
	sort.Slice(ops, func(i, j int) bool {
		return ops[i].ID.Before(ops[j].ID)
	})
}

type AppendOrderMerger struct{}

func (_ *AppendOrderMerger) Merge(operationBranches [][]timeline.Operation[Diff]) []timeline.Operation[Diff] {
	result := make([]timeline.Operation[Diff], 0)
	for _, ops := range operationBranches {
		result = append(result, ops...)
	}
	sortOperationsByID(result)
	return result
}

type ReverseAppendMerger struct{}

func (_ *ReverseAppendMerger) Merge(operationBranches [][]timeline.Operation[Diff]) []timeline.Operation[Diff] {
	result := make([]timeline.Operation[Diff], 0)
	for i := len(operationBranches) - 1; i >= 0; i-- {
		result = append(result, operationBranches[i]...)
	}
	return result
}

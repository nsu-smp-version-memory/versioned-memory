package stack

import (
	"sort"
	"sync"

	"github.com/nsu-smp-version-memory/versioned-memory/v2/internal/core"
	"github.com/nsu-smp-version-memory/versioned-memory/v2/internal/timeline"
)

type Stack struct {
	mutex           sync.Mutex
	timeline        *timeline.Timeline[Diff]
	pendingBranches []pendingBranch
	wg              sync.WaitGroup
	merger          timeline.Merger[Diff]
}

func NewStack() *Stack {
	return &Stack{
		timeline: timeline.NewTimeline[Diff](core.NewSource()),
		merger:   &TopWinsMerger{},
	}
}

func (s *Stack) Push(value int) {
	s.mutex.Lock()
	s.timeline = s.timeline.NextChange(Diff{Kind: Push, Value: value})
	s.mutex.Unlock()
}

func (s *Stack) Pop() {
	s.mutex.Lock()
	s.timeline = s.timeline.NextChange(Diff{Kind: Pop})
	s.mutex.Unlock()
}

func (s *Stack) Top() (int, bool) {
	s.mutex.Lock()
	tl := s.timeline
	s.mutex.Unlock()

	data := replayToSlice(tl)
	if len(data) == 0 {
		return 0, false
	}
	return data[len(data)-1], true
}

func (s *Stack) Items() []int {
	s.mutex.Lock()
	tl := s.timeline
	s.mutex.Unlock()

	data := replayToSlice(tl)

	out := make([]int, len(data))
	copy(out, data)
	return out
}

func (s *Stack) Size() int {
	s.mutex.Lock()
	tl := s.timeline
	s.mutex.Unlock()

	return len(replayToSlice(tl))
}

func (s *Stack) SetMerger(merger timeline.Merger[Diff]) {
	s.mutex.Lock()
	s.merger = merger
	s.mutex.Unlock()
}

func Merge(a, b *Stack) *Stack {
	merger := a.merger
	if merger == nil {
		merger = &TopWinsMerger{}
	}

	a.mutex.Lock()
	operationsA := a.timeline.Operations()
	a.mutex.Unlock()

	b.mutex.Lock()
	operationsB := b.timeline.Operations()
	b.mutex.Unlock()

	result := merger.Merge([][]timeline.Operation[Diff]{operationsA, operationsB})

	sortOperationsByID(result)

	return &Stack{
		timeline: timeline.FromOperations(core.NewSource(), result),
		merger:   merger,
	}
}

func sortOperationsByID[DIFF any](ops []timeline.Operation[DIFF]) {
	sort.Slice(ops, func(i, j int) bool {
		return ops[i].ID.Before(ops[j].ID)
	})
}

type TopWinsMerger struct {
}

func (_ *TopWinsMerger) Merge(operationBranches [][]timeline.Operation[Diff]) []timeline.Operation[Diff] {
	result := make([]timeline.Operation[Diff], 0)

	for _, ops := range operationBranches {
		result = append(result, ops...)
	}

	sortOperationsByID(result)

	return result
}

type BottomWinsMerger struct {
}

func (_ *BottomWinsMerger) Merge(operationBranches [][]timeline.Operation[Diff]) []timeline.Operation[Diff] {
	result := make([]timeline.Operation[Diff], 0)
	for i := len(operationBranches) - 1; i >= 0; i-- {
		result = append(result, operationBranches[i]...)
	}
	return result
}

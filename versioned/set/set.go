package set

import (
	"sort"
	"sync"

	"github.com/nsu-smp-version-memory/versioned-memory/internal/core"
	"github.com/nsu-smp-version-memory/versioned-memory/internal/timeline"
)

type Set struct {
	mutex           sync.Mutex
	timeline        *timeline.Timeline[Diff]
	pendingBranches []pendingBranch
	wg              sync.WaitGroup
	merger          timeline.Merger[Diff]
}

func NewSet() *Set {
	return &Set{
		timeline: timeline.NewTimeline[Diff](core.NewSource()),
		merger:   &NaturalOrderMerger{},
	}
}

func (s *Set) Add(value int) {
	s.mutex.Lock()
	s.timeline = s.timeline.NextChange(Diff{Kind: Add, Value: value})
	s.mutex.Unlock()
}

func (s *Set) Remove(value int) {
	s.mutex.Lock()
	s.timeline = s.timeline.NextChange(Diff{Kind: Remove, Value: value})
	s.mutex.Unlock()
}

func (s *Set) Contains(key int) bool {
	s.mutex.Lock()
	tl := s.timeline
	s.mutex.Unlock()

	m := replayToMap(tl)
	_, ok := m[key]
	return ok
}

func (s *Set) Items() []int {
	s.mutex.Lock()
	tl := s.timeline
	s.mutex.Unlock()

	m := replayToMap(tl)
	return mapKeysSorted(m)
}

func (s *Set) Size() int {
	s.mutex.Lock()
	tl := s.timeline
	s.mutex.Unlock()

	m := replayToMap(tl)
	return len(m)
}

func (s *Set) SetMerger(merger timeline.Merger[Diff]) {
	s.mutex.Lock()
	s.merger = merger
	s.mutex.Unlock()
}

func Merge(a, b *Set) *Set {
	merger := a.merger

	a.mutex.Lock()
	operationsA := a.timeline.Operations()
	a.mutex.Unlock()

	b.mutex.Lock()
	operationsB := b.timeline.Operations()
	b.mutex.Unlock()

	result := merger.Merge([][]timeline.Operation[Diff]{operationsA, operationsB})

	return &Set{
		timeline: timeline.FromOperations(core.NewSource(), result),
		merger:   merger,
	}
}

func sortOperationsByID[DIFF any](ops []timeline.Operation[DIFF]) {
	sort.Slice(ops, func(i, j int) bool {
		return ops[i].ID.Before(ops[j].ID)
	})
}

type NaturalOrderMerger struct {
}

func (_ *NaturalOrderMerger) Merge(operationBranches [][]timeline.Operation[Diff]) []timeline.Operation[Diff] {
	result := make([]timeline.Operation[Diff], 0)

	for _, ops := range operationBranches {
		result = append(result, ops...)
	}

	sortOperationsByID(result)

	return result
}

type ReverseOrderMerger struct {
}

func (_ *ReverseOrderMerger) Merge(operationBranches [][]timeline.Operation[Diff]) []timeline.Operation[Diff] {
	result := make([]timeline.Operation[Diff], 0)

	for i := len(operationBranches) - 1; i >= 0; i-- {
		result = append(result, operationBranches[i]...)
	}

	return result
}

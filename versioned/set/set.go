package set

import (
	"sort"
	"sync"

	"github.com/nsu-smp-version-memory/versioned-memory/internal/core"
)

type Set struct {
	mutex           sync.Mutex
	timeline        *core.Timeline[Diff]
	pendingBranches []pendingBranch
}

func NewSet() *Set {
	return &Set{
		timeline: core.NewTimeline[Diff](core.NewSource()),
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

func Merge(a, b *Set) *Set {
	operations := make([]core.Operation[Diff], 0)

	a.mutex.Lock()
	timelineA := a.timeline
	a.mutex.Unlock()
	operations = append(operations, timelineA.Operations()...)

	b.mutex.Lock()
	timelineB := b.timeline
	b.mutex.Unlock()
	operations = append(operations, timelineB.Operations()...)

	sortOperationsByID(operations)

	return &Set{
		timeline: core.TimelineFromOperations(core.NewSource(), operations),
	}
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

func sortOperationsByID[DIFF any](ops []core.Operation[DIFF]) {
	sort.Slice(ops, func(i, j int) bool {
		return ops[i].ID.Before(ops[j].ID)
	})
}

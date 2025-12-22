package set

import (
	"sort"

	"github.com/nsu-smp-version-memory/versioned-memory/internal/core"
)

type Set struct {
	timeline *core.Timeline[Diff]
}

func NewSet() *Set {
	return &Set{
		timeline: core.NewTimeline[Diff](core.NewSource()),
	}
}

func (s *Set) Add(value int) {
	s.timeline = s.timeline.NextChange(Diff{Kind: Add, Value: value})
}

func (s *Set) Remove(value int) {
	s.timeline = s.timeline.NextChange(Diff{Kind: Remove, Value: value})
}

func Merge(a, b *Set) *Set {
	ops := make([]core.Operation[Diff], 0)

	ops = append(ops, a.timeline.Operations()...)
	ops = append(ops, b.timeline.Operations()...)

	sort.Slice(ops, func(i, j int) bool { return ops[i].ID.Before(ops[j].ID) })

	return &Set{
		timeline: core.TimelineFromOperations(core.NewSource(), ops),
	}
}

func (s *Set) Contains(key int) bool {
	m := replayToMap(s.timeline)
	_, ok := m[key]
	return ok
}

func (s *Set) Items() []int {
	m := replayToMap(s.timeline)
	return mapKeysSorted(m)
}

func (s *Set) Size() int {
	m := replayToMap(s.timeline)
	return len(m)
}

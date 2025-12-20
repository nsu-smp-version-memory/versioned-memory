package versioned

import (
	"maps"

	"github.com/nsu-smp-version-memory/versioned-memory/internal/timeline"
)

type setDiff struct {
	added   map[int]bool
	removed map[int]bool
}

func makeSetDiff() setDiff {
	return setDiff{
		added:   make(map[int]bool),
		removed: make(map[int]bool),
	}
}

type Set struct {
	t      *timeline.Timeline[setDiff]
	values map[int]bool
}

func NewSet() *Set {
	return &Set{
		t:      timeline.New[setDiff](),
		values: make(map[int]bool),
	}
}

func (s *Set) Add(v int) *Set {
	diff := makeSetDiff()
	diff.added[v] = true
	values := maps.Clone(s.values)
	values[v] = true

	return &Set{
		t:      s.t.NextChange(diff),
		values: values,
	}
}

func (s *Set) Remove(v int) *Set {
	diff := makeSetDiff()
	diff.removed[v] = true
	values := maps.Clone(s.values)
	delete(values, v)

	return &Set{
		t:      s.t.NextChange(diff),
		values: values,
	}
}

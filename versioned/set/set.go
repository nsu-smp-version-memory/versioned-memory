package set

import (
	"sort"

	"github.com/nsu-smp-version-memory/versioned-memory/internal/core"
)

type Set struct {
	src *core.Source
	log []operation
}

func NewSet() *Set {
	return &Set{
		src: core.NewSourceFor(core.KindSet),
		log: nil,
	}
}

func (s *Set) Add(value int) {
	s.log = append(s.log, operation{
		id:    s.src.NextOperationID(),
		kind:  operationAdd,
		value: value,
	})
}

func (s *Set) Remove(value int) {
	s.log = append(s.log, operation{
		id:    s.src.NextOperationID(),
		kind:  operationRemove,
		value: value,
	})
}

func Merge(a, b *Set) *Set {
	out := NewSet()
	out.log = append(out.log, a.log...)
	out.log = append(out.log, b.log...)

	sort.Slice(out.log, func(i, j int) bool { return out.log[i].id.Before(out.log[j].id) })

	return out
}

func (s *Set) Contains(key int) bool {
	m := replayToMap(s.log)
	_, ok := m[key]
	return ok
}

func (s *Set) Items() []int {
	m := replayToMap(s.log)
	return mapKeysSorted(m)
}

func (s *Set) Size() int {
	m := replayToMap(s.log)
	return len(m)
}

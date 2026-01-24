package timeline

import (
	"slices"

	"github.com/nsu-smp-version-memory/versioned-memory/internal/core"
)

type Operation[DIFF any] struct {
	ID   core.OperationID
	Diff DIFF
}

type Node[DIFF any] struct {
	id   core.OperationID
	prev *Node[DIFF]
	diff DIFF
}

type Timeline[DIFF any] struct {
	last   *Node[DIFF]
	source *core.Source
}

func NewTimeline[DIFF any](src *core.Source) *Timeline[DIFF] {
	return &Timeline[DIFF]{last: nil, source: src}
}

func (t *Timeline[DIFF]) NextChange(diff DIFF) *Timeline[DIFF] {
	node := &Node[DIFF]{
		id:   t.source.NextOperationID(),
		prev: t.last,
		diff: diff,
	}

	return &Timeline[DIFF]{
		last:   node,
		source: t.source,
	}
}

func FromOperations[DIFF any](src *core.Source, ops []Operation[DIFF]) *Timeline[DIFF] {
	var last *Node[DIFF]
	for i := 0; i < len(ops); i++ {
		last = &Node[DIFF]{
			id:   ops[i].ID,
			prev: last,
			diff: ops[i].Diff,
		}
	}
	return &Timeline[DIFF]{last: last, source: src}
}

func (t *Timeline[DIFF]) Operations() []Operation[DIFF] {
	tmp := make([]Operation[DIFF], 0)
	for cur := t.last; cur != nil; cur = cur.prev {
		tmp = append(tmp, Operation[DIFF]{ID: cur.id, Diff: cur.diff})
	}
	slices.Reverse(tmp)
	return tmp
}

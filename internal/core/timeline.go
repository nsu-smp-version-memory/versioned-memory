package core

type Operation[DIFF any] struct {
	ID   OperationID
	Diff DIFF
}

type Node[DIFF any] struct {
	id   OperationID
	prev *Node[DIFF]
	diff DIFF
}

type Timeline[DIFF any] struct {
	last   *Node[DIFF]
	source *Source
}

func NewTimeline[DIFF any](src *Source) *Timeline[DIFF] {
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

func TimelineFromOperations[DIFF any](src *Source, ops []Operation[DIFF]) *Timeline[DIFF] {
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
	var tmp []Operation[DIFF]
	for cur := t.last; cur != nil; cur = cur.prev {
		tmp = append(tmp, Operation[DIFF]{ID: cur.id, Diff: cur.diff})
	}
	for i, j := 0, len(tmp)-1; i < j; i, j = i+1, j-1 {
		tmp[i], tmp[j] = tmp[j], tmp[i]
	}
	return tmp
}

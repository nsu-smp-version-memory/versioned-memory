package core

type ForkPoint[DIFF any] struct {
	node *Node[DIFF]
}

func (t *Timeline[DIFF]) Fork(newSource *Source) (*Timeline[DIFF], ForkPoint[DIFF]) {
	return &Timeline[DIFF]{
		last:   t.last,
		source: newSource,
	}, ForkPoint[DIFF]{node: t.last}
}

func (t *Timeline[DIFF]) OperationsAfter(p ForkPoint[DIFF]) []Operation[DIFF] {
	var tmp []Operation[DIFF]
	for cur := t.last; cur != nil && cur != p.node; cur = cur.prev {
		tmp = append(tmp, Operation[DIFF]{ID: cur.id, Diff: cur.diff})
	}

	for i, j := 0, len(tmp)-1; i < j; i, j = i+1, j-1 {
		tmp[i], tmp[j] = tmp[j], tmp[i]
	}
	return tmp
}

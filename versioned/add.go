package versioned

import "github.com/nsu-smp-version-memory/versioned-memory/internal/core"

func (s Set) Add(src *core.Source, key int) Set {
	operationID := src.NextOperationID()

	var changed bool
	newRoot := addNode(s.root, key, operationID, &changed)
	if !changed {
		return s
	}

	return Set{
		root:    newRoot,
		version: src.Version(),
	}
}

func addNode(n *node, value int, op core.OperationID, changed *bool) *node {
	if n == nil {
		*changed = true
		return &node{
			value: value,
			adds:  TagSet{op: {}},
		}
	}

	if value < n.value {
		left := addNode(n.left, value, op, changed)
		if !*changed {
			return n
		}
		return &node{
			value: n.value,
			adds:  n.adds,
			dels:  n.dels,
			left:  left,
			right: n.right,
		}
	}

	if value > n.value {
		right := addNode(n.right, value, op, changed)
		if !*changed {
			return n
		}
		return &node{
			value: n.value,
			adds:  n.adds,
			dels:  n.dels,
			left:  n.left,
			right: right,
		}
	}

	if n.adds != nil {
		if _, exists := n.adds[op]; exists {
			*changed = false
			return n
		}
	}

	adds := n.adds.Copy()
	if adds == nil {
		adds = make(TagSet)
	}
	adds[op] = struct{}{}

	*changed = true
	return &node{
		value: n.value,
		adds:  adds,
		dels:  n.dels,
		left:  n.left,
		right: n.right,
	}
}

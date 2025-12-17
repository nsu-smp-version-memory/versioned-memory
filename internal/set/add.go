package set

import "github.com/nsu-smp-version-memory/version_memory/internal/core"

func (set Set) Add(src *core.Source, key int) Set {
	op := src.NextOperationID()

	var changed bool
	newRoot := addNode(set.root, key, op, &changed)
	if !changed {
		return set
	}

	return Set{
		root:    newRoot,
		version: src.Version(),
	}
}

func addNode(n *node, key int, op core.OperationID, changed *bool) *node {
	if n == nil {
		*changed = true
		return &node{
			key:  key,
			adds: TagSet{op: {}},
		}
	}

	if key < n.key {
		left := addNode(n.left, key, op, changed)
		if !*changed {
			return n
		}
		return &node{
			key:   n.key,
			adds:  n.adds,
			dels:  n.dels,
			left:  left,
			right: n.right,
		}
	}

	if key > n.key {
		right := addNode(n.right, key, op, changed)
		if !*changed {
			return n
		}
		return &node{
			key:   n.key,
			adds:  n.adds,
			dels:  n.dels,
			left:  n.left,
			right: right,
		}
	}

	// key == n.key
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
		key:   n.key,
		adds:  adds,
		dels:  n.dels,
		left:  n.left,
		right: n.right,
	}
}

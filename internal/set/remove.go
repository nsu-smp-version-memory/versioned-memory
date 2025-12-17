package set

import "github.com/nsu-smp-version-memory/version_memory/internal/core"

func (set Set) Remove(src *core.Source, value int) Set {
	_ = src.NextOperationID()

	var changed bool
	newRoot := removeNode(s.root, value, &changed)
	if !changed {
		return set
	}

	return Set{
		root:    newRoot,
		version: src.Version(),
	}
}

func removeNode(n *node, value int, changed *bool) *node {
	if n == nil {
		return nil
	}

	if value < n.value {
		left := removeNode(n.left, value, changed)
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
		right := removeNode(n.right, value, changed)
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

	if len(n.adds) == 0 {
		return n
	}

	dels := n.dels.Copy()
	if dels == nil {
		dels = make(TagSet)
	}

	localChanged := false
	for tag := range n.adds {
		if _, already := dels[tag]; already {
			continue
		}
		dels[tag] = struct{}{}
		localChanged = true
	}

	if !localChanged {
		return n
	}

	*changed = true
	return &node{
		value: n.value,
		adds:  n.adds,
		dels:  dels,
		left:  n.left,
		right: n.right,
	}
}

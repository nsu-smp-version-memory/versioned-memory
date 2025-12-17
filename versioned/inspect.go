package versioned

func (s Set) Keys() []int {
	var out []int
	collectAliveInOrder(s.root, &out)
	return out
}

func collectAliveInOrder(n *node, out *[]int) {
	if n == nil {
		return
	}
	collectAliveInOrder(n.left, out)

	alive := false
	for tag := range n.adds {
		if _, removed := n.dels[tag]; !removed {
			alive = true
			break
		}
	}
	if alive {
		*out = append(*out, n.value)
	}

	collectAliveInOrder(n.right, out)
}

func (s Set) Height() int {
	return height(s.root)
}

func height(n *node) int {
	if n == nil {
		return 0
	}
	hl := height(n.left)
	hr := height(n.right)
	if hl > hr {
		return hl + 1
	}
	return hr + 1
}

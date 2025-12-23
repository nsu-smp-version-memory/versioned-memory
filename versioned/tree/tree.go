package tree

import (
	"sort"
	"sync"

	"github.com/nsu-smp-version-memory/versioned-memory/internal/core"
)

type node struct {
	left   *node
	right  *node
	value  int
	height int
}

func newNode(value int) *node {
	return &node{
		left:   nil,
		right:  nil,
		value:  value,
		height: 1,
	}
}

func safeHeight(n *node) int {
	if n == nil {
		return 0
	} else {
		return n.height
	}
}

func (n *node) balanceFactor() int {
	return safeHeight(n.right) - safeHeight(n.left)
}

func (n *node) fixHeight() *node {
	leftHeight := safeHeight(n.left)
	rightHeight := safeHeight(n.right)

	if leftHeight > rightHeight {
		return &node{
			left:   n.left,
			right:  n.right,
			value:  n.value,
			height: leftHeight + 1,
		}
	} else {
		return &node{
			left:   n.left,
			right:  n.right,
			value:  n.value,
			height: rightHeight + 1,
		}
	}
}

func (n *node) setLeft(l *node) *node {
	return &node{
		left:   l,
		right:  n.right,
		value:  n.value,
		height: n.height,
	}
}

func (n *node) setRight(r *node) *node {
	return &node{
		left:   n.left,
		right:  r,
		value:  n.value,
		height: n.height,
	}
}

func (n *node) rotateRight() *node {
	p := n
	q := n.left

	p = p.setLeft(q.right)
	p = p.fixHeight()

	q = q.setRight(p)
	q = q.fixHeight()

	return q
}

func (n *node) rotateLeft() *node {
	p := n.right
	q := n

	q = q.setRight(p.left)
	q = q.fixHeight()

	p = p.setLeft(q)
	p = p.fixHeight()

	return p
}

func (n *node) balance() *node {
	p := n

	p = p.fixHeight()
	if p.balanceFactor() == 2 {
		if p.right.balanceFactor() < 0 {
			p = p.setRight(p.right.rotateRight())
		}

		return p.rotateLeft()
	}

	if p.balanceFactor() == -2 {
		if p.left.balanceFactor() > 0 {
			p = p.setLeft(n.left.rotateLeft())
		}

		return p.rotateRight()
	}

	return p
}

func (n *node) insert(k int) *node {
	p := n

	if k < p.value {
		if p.left == nil {
			p = p.setLeft(newNode(k))
		} else {
			p = p.setLeft(p.left.insert(k))
		}
	} else {
		if p.right == nil {
			p = p.setRight(newNode(k))
		} else {
			p = p.setRight(p.right.insert(k))
		}
	}

	return p.balance()
}

func (n *node) findMin() *node {
	if n.left == nil {
		return n
	} else {
		return n.left.findMin()
	}
}

func (n *node) removeMin() *node {
	p := n

	if p.left == nil {
		return p.right
	}

	p = p.setLeft(p.left.removeMin())
	return p.balance()
}

func (n *node) remove(value int) *node {
	p := n

	if value < p.value {
		if p.left == nil {
			return p
		} else {
			p = p.setLeft(p.left.remove(value))
			return p.balance()
		}
	} else if value > n.value {
		if p.right == nil {
			return p
		} else {
			p = p.setRight(p.right.remove(value))
			return p.balance()
		}
	} else {
		q := p.left
		r := p.right
		if r == nil {
			return q
		}
		m := r.findMin()
		m = m.setRight(r.removeMin())
		m = m.setLeft(q)
		return m.balance()
	}
}

func (n *node) contains(value int) bool {
	if value < n.value {
		if n.left == nil {
			return false
		} else {
			return n.left.contains(value)
		}
	} else if value > n.value {
		if n.right == nil {
			return false
		} else {
			return n.right.contains(value)
		}
	} else {
		return true
	}
}

func (n *node) items() []int {
	result := make([]int, 0)
	if n.left != nil {
		result = append(result, n.left.items()...)
	}

	result = append(result, n.value)

	if n.right != nil {
		result = append(result, n.right.items()...)
	}

	return result
}

func (n *node) size() int {
	result := 1

	if n.left != nil {
		result += n.left.size()
	}

	if n.right != nil {
		result += n.right.size()
	}

	return result
}

type Tree struct {
	mutex           sync.Mutex
	timeline        *core.Timeline[Diff]
	pendingBranches []pendingBranch
	wg              sync.WaitGroup
	merger          core.Merger[Diff]
	root            *node
}

func NewTree() *Tree {
	return &Tree{
		timeline: core.NewTimeline[Diff](core.NewSource()),
		merger:   &NaturalOrderMerger{},
	}
}

func (t *Tree) Add(value int) {
	t.mutex.Lock()
	t.timeline = t.timeline.NextChange(Diff{Kind: Add, Value: value})
	if t.root == nil {
		t.root = newNode(value)
	} else {
		t.root = t.root.insert(value)
	}
	t.mutex.Unlock()
}

func (t *Tree) Remove(value int) {
	t.mutex.Lock()
	t.timeline = t.timeline.NextChange(Diff{Kind: Remove, Value: value})
	if t.root != nil {
		t.root = t.root.remove(value)
	}
	t.mutex.Unlock()
}

func (t *Tree) Contains(key int) bool {
	t.mutex.Lock()
	root := t.root
	t.mutex.Unlock()

	if root == nil {
		return false
	} else {
		return root.contains(key)
	}
}

func (t *Tree) Items() []int {
	t.mutex.Lock()
	root := t.root
	t.mutex.Unlock()

	if root == nil {
		return make([]int, 0)
	} else {
		return root.items()
	}
}

func (t *Tree) Size() int {
	t.mutex.Lock()
	root := t.root
	t.mutex.Unlock()

	if root == nil {
		return 0
	} else {
		return root.size()
	}
}

func (t *Tree) SetMerger(merger core.Merger[Diff]) {
	t.mutex.Lock()
	t.merger = merger
	t.mutex.Unlock()
}

func Merge(a, b *Tree) *Tree {
	merger := a.merger

	a.mutex.Lock()
	operationsA := a.timeline.Operations()
	root := a.root
	a.mutex.Unlock()

	b.mutex.Lock()
	operationsB := b.timeline.Operations()
	b.mutex.Unlock()

	for i := len(operationsA) - 1; i >= 0; i-- {
		switch operationsA[i].Diff.Kind {
		case Add:
			if root != nil {
				root = root.remove(operationsA[i].Diff.Value)
			}
		case Remove:
			if root == nil {
				root = newNode(operationsA[i].Diff.Value)
			} else {
				root = root.insert(operationsA[i].Diff.Value)
			}
		}
	}

	result := merger.Merge([][]core.Operation[Diff]{operationsA, operationsB})

	for _, op := range result {
		switch op.Diff.Kind {
		case Add:
			if root == nil {
				root = newNode(op.Diff.Value)
			} else {
				root = root.insert(op.Diff.Value)
			}
		case Remove:
			if root != nil {
				root = root.remove(op.Diff.Value)
			}
		}
	}

	return &Tree{
		timeline: core.TimelineFromOperations(core.NewSource(), result),
		merger:   merger,
		root:     root,
	}
}

func sortOperationsByID[DIFF any](ops []core.Operation[DIFF]) {
	sort.Slice(ops, func(i, j int) bool {
		return ops[i].ID.Before(ops[j].ID)
	})
}

type NaturalOrderMerger struct {
}

func (_ *NaturalOrderMerger) Merge(operationBranches [][]core.Operation[Diff]) []core.Operation[Diff] {
	result := make([]core.Operation[Diff], 0)

	for _, ops := range operationBranches {
		result = append(result, ops...)
	}

	sortOperationsByID(result)

	return result
}

type ReverseOrderMerger struct {
}

func (_ *ReverseOrderMerger) Merge(operationBranches [][]core.Operation[Diff]) []core.Operation[Diff] {
	result := make([]core.Operation[Diff], 0)

	for i := len(operationBranches) - 1; i >= 0; i-- {
		result = append(result, operationBranches[i]...)
	}

	return result
}

package set

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicAddRemoveContainsItems(t *testing.T) {
	s := NewSet()

	s.Add(1)
	s.Add(2)
	s.Add(3)
	s.Remove(2)

	assert.True(t, s.Contains(1))
	assert.False(t, s.Contains(2))
	assert.True(t, s.Contains(3))

	assert.Equal(t, 2, s.Size())

	assert.Contains(t, s.Items(), 1)
	assert.Contains(t, s.Items(), 3)
}

func TestMergeNonConflicting(t *testing.T) {
	a := NewSet()
	b := NewSet()

	a.Add(1)
	a.Add(2)

	b.Add(3)
	b.Remove(2)

	m := Merge(a, b)

	assert.Contains(t, m.Items(), 1)
	assert.Contains(t, m.Items(), 3)
	assert.Equal(t, 2, m.Size())
}

func TestMergeConflictingSameKeyDeterministicBySourceOrder(t *testing.T) {
	{
		a := NewSet()
		b := NewSet()

		a.Add(1)
		b.Remove(1)

		m := Merge(a, b)
		assert.Empty(t, m.Items())
	}

	{
		b := NewSet()
		a := NewSet()

		a.Add(1)
		b.Remove(1)

		m := Merge(a, b)

		assert.Equal(t, 1, m.Size())
		assert.Contains(t, m.Items(), 1)
	}
}

func TestMergeAssociative(t *testing.T) {
	a := NewSet()
	b := NewSet()
	c := NewSet()

	a.Add(1)
	a.Add(2)

	b.Add(3)
	b.Remove(2)

	c.Add(4)
	c.Remove(1)

	left := Merge(Merge(a, b), c)
	right := Merge(a, Merge(b, c))

	assert.True(t, reflect.DeepEqual(left.Items(), right.Items()))
}

func TestMergeDoesNotMutateInputs(t *testing.T) {
	a := NewSet()
	b := NewSet()

	a.Add(1)
	b.Add(2)

	beforeA := a.Items()
	beforeB := b.Items()

	_ = Merge(a, b)

	afterA := a.Items()
	afterB := b.Items()

	assert.True(t, reflect.DeepEqual(afterA, beforeA))
	assert.True(t, reflect.DeepEqual(afterB, beforeB))
}

func TestConcurrentTwoGoroutinesMergeAdds(t *testing.T) {
	doneA := make(chan *Set)
	doneB := make(chan *Set)

	go func() {
		a := NewSet()
		a.Add(1)
		a.Add(2)
		doneA <- a
	}()

	go func() {
		b := NewSet()
		b.Add(3)
		b.Add(4)
		doneB <- b
	}()

	a := <-doneA
	b := <-doneB

	m := Merge(a, b)

	got := m.Items()
	want := []int{1, 2, 3, 4}

	assert.True(t, reflect.DeepEqual(got, want))
}
func TestWithBranchAndMergeBranches(t *testing.T) {
	s := NewSet()
	s.Add(1)
	s.Add(2)

	done := s.WithBranch(func(s *Set) {
		s.Remove(1)
		s.Add(10)
	})

	s.Add(3)

	<-done
	s.MergeBranches()

	assert.True(t, reflect.DeepEqual(s.Items(), []int{2, 3, 10}))
}

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

	if got, want := m.Items(), []int{1, 3}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Merged Items()=%v, want %v", got, want)
	}
}

func TestMergeConflictingSameKeyDeterministicBySourceOrder(t *testing.T) {
	{
		a := NewSet()
		b := NewSet()

		a.Add(1)
		b.Remove(1)

		m := Merge(a, b)
		if got, want := m.Items(), []int{}; !reflect.DeepEqual(got, want) {
			t.Fatalf("[A then B] Items()=%v, want %v", got, want)
		}
	}

	{
		b := NewSet()
		a := NewSet()

		a.Add(1)
		b.Remove(1)

		m := Merge(a, b)
		if got, want := m.Items(), []int{1}; !reflect.DeepEqual(got, want) {
			t.Fatalf("[B then A] Items()=%v, want %v", got, want)
		}
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

	if gotL, gotR := left.Items(), right.Items(); !reflect.DeepEqual(gotL, gotR) {
		t.Fatalf("merge not associative: left=%v right=%v", gotL, gotR)
	}
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

	if !reflect.DeepEqual(beforeA, afterA) {
		t.Fatalf("input A mutated: before=%v after=%v", beforeA, afterA)
	}
	if !reflect.DeepEqual(beforeB, afterB) {
		t.Fatalf("input B mutated: before=%v after=%v", beforeB, afterB)
	}
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
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("merged items=%v, want %v", got, want)
	}
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

	if got, want := s.Items(), []int{2, 3, 10}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Items()=%v, want %v", got, want)
	}
}

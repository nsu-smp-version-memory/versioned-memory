package stack

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicPushPopTopItems(t *testing.T) {
	s := NewStack()

	s.Push(1)
	s.Push(2)
	s.Push(3)

	v, ok := s.Top()
	assert.True(t, ok)
	assert.Equal(t, 3, v)

	s.Pop()
	v, ok = s.Top()
	assert.True(t, ok)
	assert.Equal(t, 2, v)

	assert.True(t, reflect.DeepEqual(s.Items(), []int{1, 2}))
	assert.Equal(t, 2, s.Size())
}

func TestPopOnEmptyIsNoOp(t *testing.T) {
	s := NewStack()
	s.Pop()
	_, ok := s.Top()
	assert.False(t, ok)
	assert.Equal(t, 0, s.Size())
}

func TestMergeDeterministic(t *testing.T) {
	a := NewStack()
	b := NewStack()

	a.Push(1)
	a.Push(2)
	b.Push(3)

	m := Merge(a, b)
	assert.Equal(t, 3, m.Size())
	assert.True(t, reflect.DeepEqual(m.Items(), []int{1, 2, 3}) || reflect.DeepEqual(m.Items(), []int{3, 1, 2}))
}

func TestWithBranchAndMergeBranches(t *testing.T) {
	s := NewStack()
	s.Push(1)
	s.Push(2)

	done := s.WithBranch(func(local *Stack) {
		local.Pop()
		local.Push(10)
	})

	s.Push(3)

	<-done
	s.MergeBranches()

	assert.True(t, reflect.DeepEqual(s.Items(), []int{1, 2, 10}))
}

func TestGoWithJoin(t *testing.T) {
	s := NewStack()
	s.Push(1)
	s.Push(2)

	s.Go(func(local *Stack) {
		local.Pop()
		local.Push(10)
	})

	s.Push(3)

	s.Join()
	assert.True(t, reflect.DeepEqual(s.Items(), []int{1, 2, 10}))
}

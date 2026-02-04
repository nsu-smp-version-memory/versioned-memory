package queue

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicEnqueueDequeueFrontItems(t *testing.T) {
	q := NewQueue()

	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	v, ok := q.Front()
	assert.True(t, ok)
	assert.Equal(t, 1, v)

	q.Dequeue()
	v, ok = q.Front()
	assert.True(t, ok)
	assert.Equal(t, 2, v)

	assert.True(t, reflect.DeepEqual(q.Items(), []int{2, 3}))
	assert.Equal(t, 2, q.Size())
}

func TestDequeueOnEmptyIsNoOp(t *testing.T) {
	q := NewQueue()
	q.Dequeue()
	_, ok := q.Front()
	assert.False(t, ok)
	assert.Equal(t, 0, q.Size())
}

func TestWithBranchAndMergeBranches(t *testing.T) {
	q := NewQueue()
	q.Enqueue(1)
	q.Enqueue(2)

	done := q.WithBranch(func(local *Queue) {
		local.Dequeue()
		local.Enqueue(10)
	})

	q.Enqueue(3)

	<-done
	q.MergeBranches()

	assert.True(t, reflect.DeepEqual(q.Items(), []int{2, 3, 10}))
}

func TestGoWithJoin(t *testing.T) {
	q := NewQueue()
	q.Enqueue(1)
	q.Enqueue(2)

	q.Go(func(local *Queue) {
		local.Dequeue()
		local.Enqueue(10)
	})

	q.Enqueue(3)

	q.Join()
	assert.True(t, reflect.DeepEqual(q.Items(), []int{2, 3, 10}))
}

func TestMergeDoesNotMutateInputs(t *testing.T) {
	a := NewQueue()
	b := NewQueue()

	a.Enqueue(1)
	b.Enqueue(2)

	beforeA := a.Items()
	beforeB := b.Items()

	_ = Merge(a, b)

	assert.True(t, reflect.DeepEqual(beforeA, a.Items()))
	assert.True(t, reflect.DeepEqual(beforeB, b.Items()))
}

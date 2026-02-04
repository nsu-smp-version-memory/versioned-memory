package queue

type Kind uint8

const (
	Enqueue Kind = iota + 1
	Dequeue
)

type Diff struct {
	Kind  Kind
	Value int
}

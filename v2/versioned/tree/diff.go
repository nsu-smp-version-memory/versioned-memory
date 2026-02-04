package tree

type Kind uint8

const (
	Add Kind = iota + 1
	Remove
)

type Diff struct {
	Kind  Kind
	Value int
}

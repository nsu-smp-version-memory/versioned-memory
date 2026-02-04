package stack

type Kind uint8

const (
	Push Kind = iota + 1
	Pop
)

type Diff struct {
	Kind  Kind
	Value int
}

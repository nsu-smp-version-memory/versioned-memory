package core

type Kind uint8

const (
	KindSet Kind = iota + 1
	KindStack
	KindQueue
)

func (k Kind) String() string {
	switch k {
	case KindSet:
		return "set"
	case KindStack:
		return "stack"
	case KindQueue:
		return "queue"
	default:
		return "unknown"
	}
}

package set

import "github.com/nsu-smp-version-memory/versioned-memory/internal/core"

type operationKind uint8

const (
	operationAdd operationKind = iota + 1
	operationRemove
)

type operation struct {
	id    core.OperationID
	kind  operationKind
	value int
}

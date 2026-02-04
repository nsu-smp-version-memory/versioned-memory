package core

import (
	"encoding/binary"
	"fmt"
)

type SourceID uint32
type OperationIndex uint32

type OperationID uint64

func NewOperationID(source SourceID, index OperationIndex) OperationID {
	return OperationID(uint64(source)<<32 | uint64(index))
}

func (id OperationID) Source() SourceID {
	return SourceID(uint64(id) >> 32)
}

func (id OperationID) Index() OperationIndex {
	return OperationIndex(uint64(id) & 0xffffffff)
}

func (id OperationID) Before(other OperationID) bool {
	return uint64(id) < uint64(other)
}

func (id OperationID) String() string {
	return fmt.Sprintf("%d:%d", id.Source(), id.Index())
}

func (id OperationID) MarshalBinary() ([]byte, error) {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(id))
	return b[:], nil
}

func (id *OperationID) UnmarshalBinary(data []byte) error {
	if len(data) != 8 {
		return fmt.Errorf("invalid OperationID length: got %d, want 8", len(data))
	}
	*id = OperationID(binary.BigEndian.Uint64(data))
	return nil
}

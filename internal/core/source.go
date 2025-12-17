package core

import (
	"sync"
	"sync/atomic"
)

type Source struct {
	id SourceID

	next  atomic.Uint32
	armed atomic.Bool

	versionManager *VersionManager

	mutex   sync.Mutex
	current *Version
}

func NewSource(id SourceID, versionManager *VersionManager) *Source {
	s := &Source{
		id:             id,
		versionManager: versionManager,
	}
	s.current = versionManager.Root()
	s.armed.Store(true)
	return s
}

func (s *Source) ID() SourceID {
	return s.id
}

func (s *Source) Version() *Version {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.current
}

func (s *Source) Arm() {
	s.armed.Store(true)
}

func (source *Source) NextOperationID() OperationID {
	if source.armed.Swap(false) {
		source.mutex.Lock()
		source.current = source.versionManager.NewChild(source.current)
		source.mutex.Unlock()
	}

	idx := OperationIndex(s.next.Add(1))
	return NewOperationID(s.id, idx)
}

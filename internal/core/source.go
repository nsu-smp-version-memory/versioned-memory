package core

import (
	"sync"
	"sync/atomic"
)

type Source struct {
	id SourceID

	next  atomic.Uint32
	armed atomic.Bool

	version_manager *VersionManager

	mutex   sync.Mutex
	current *Version
}

func NewSource(id SourceID, version_manager *VersionManager) *Source {
	s := &Source{
		id:              id,
		version_manager: version_manager,
	}
	s.current = version_manager.Root()
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

func (s *Source) NextOperationID() OperationID {
	if s.armed.Swap(false) {
		s.mutex.Lock()
		s.current = s.version_manager.NewChild(s.current)
		s.mutex.Unlock()
	}

	idx := OperationIndex(s.next.Add(1))
	return NewOperationID(s.id, idx)
}

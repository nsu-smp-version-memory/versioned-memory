package core

type Source struct {
	id SourceID

	next  OperationIndex
	armed bool

	versionManager *VersionManager
	current        *Version
}

func NewSource(id SourceID, vm *VersionManager) *Source {
	s := &Source{
		id:             id,
		versionManager: vm,
		current:        vm.Root(),
		armed:          true,
	}
	return s
}

func (s *Source) ID() SourceID {
	return s.id
}

func (s *Source) Version() *Version {
	return s.current
}

func (s *Source) Arm() {
	s.armed = true
}

func (s *Source) NextOperationID() OperationID {
	if s.armed {
		s.current = s.versionManager.NewChild(s.current)
		s.armed = false
	}

	s.next++
	return NewOperationID(s.id, s.next)
}

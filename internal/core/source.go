package core

type Source struct {
	id SourceID

	next  OperationIndex
	armed bool

	manager *Manager
	current *Version
	kind    Kind
}

func newSource(kind Kind, m *Manager) *Source {
	s := &Source{
		id:      m.newSourceID(),
		manager: m,
		kind:    kind,
		current: m.Root(kind),
		armed:   true,
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
		s.current = s.manager.NewChild(s.current)
		s.armed = false
	}
	s.next++
	return NewOperationID(s.id, s.next)
}

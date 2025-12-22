package core

type Source struct {
	id   SourceID
	next OperationIndex
}

func NewSource() *Source {
	return &Source{
		id:   DefaultManager().NewSourceID(),
		next: 0,
	}
}

func (s *Source) ID() SourceID {
	return s.id
}

func (s *Source) NextOperationID() OperationID {
	s.next++
	return NewOperationID(s.id, s.next)
}

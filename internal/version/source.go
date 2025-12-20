package version

import "sync/atomic"

type Source struct {
	currVersion atomic.Uint64
}

func NewSource() *Source {
	return &Source{
		currVersion: atomic.Uint64{},
	}
}

func (s *Source) Peek() uint64 {
	return s.currVersion.Add(1)
}

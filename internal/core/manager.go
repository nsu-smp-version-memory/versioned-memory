package core

import "sync"

type Manager struct {
	mu         sync.Mutex
	nextSource SourceID
}

func newManager() *Manager {
	return &Manager{
		nextSource: 1,
	}
}

func (m *Manager) NewSourceID() SourceID {
	m.mu.Lock()
	defer m.mu.Unlock()
	id := m.nextSource
	m.nextSource++
	return id
}

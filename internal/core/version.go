package core

import "sync"

type VersionID uint64

type Version struct {
	id     VersionID
	parent *Version
}

func (v *Version) ID() VersionID {
	return v.id
}

func (v *Version) Parent() *Version {
	return v.parent
}

type Manager struct {
	mutex       sync.Mutex
	nextVersion VersionID
	nextSource  SourceID
	roots       map[Kind]*Version
}

func newManager() *Manager {
	manager := &Manager{
		nextVersion: 1,
		nextSource:  1,
		roots:       make(map[Kind]*Version, 3),
	}
	// One root per container kind.
	manager.roots[KindSet] = &Version{id: manager.allocVersionIDLocked()}
	manager.roots[KindStack] = &Version{id: manager.allocVersionIDLocked()}
	manager.roots[KindQueue] = &Version{id: manager.allocVersionIDLocked()}
	return manager
}

func (m *Manager) allocVersionIDLocked() VersionID {
	id := m.nextVersion
	m.nextVersion++
	return id
}

func (m *Manager) Root(kind Kind) *Version {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.roots[kind]
}

func (m *Manager) NewChild(parent *Version) *Version {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if parent == nil {
		parent = &Version{id: m.allocVersionIDLocked()}
	}

	version := &Version{
		id:     m.allocVersionIDLocked(),
		parent: parent,
	}
	return version
}

func (m *Manager) newSourceID() SourceID {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	id := m.nextSource
	m.nextSource++
	return id
}

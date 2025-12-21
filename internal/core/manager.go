package core

import "sync"

var (
	defaultOnce    sync.Once
	defaultManager *Manager
)

func Default() *Manager {
	defaultOnce.Do(func() {
		defaultManager = newManager()
	})
	return defaultManager
}

func NewSourceFor(kind Kind) *Source {
	return newSource(kind, Default())
}

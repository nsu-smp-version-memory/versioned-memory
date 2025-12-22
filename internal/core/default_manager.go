package core

import "sync"

var (
	defaultOnce    sync.Once
	defaultManager *Manager
)

func DefaultManager() *Manager {
	defaultOnce.Do(func() {
		defaultManager = newManager()
	})
	return defaultManager
}

package core

import "sync"

type VersionID uint64

type Version struct {
	id     VersionID
	parent *Version
}

func (version *Version) ID() VersionID {
	if version == nil {
		return 0
	}
	return version.id
}

func (version *Version) Parent() *Version {
	if version == nil {
		return nil
	}
	return version.parent
}

type VersionManager struct {
	mutex sync.Mutex
	next  VersionID
	root  *Version
}

func NewVersionManager() *VersionManager {
	version_manager := &VersionManager{next: 1}
	version_manager.root = &Version{id: version_manager.next}
	version_manager.next++
	return version_manager
}

func (version_manager *VersionManager) Root() *Version {
	version_manager.mutex.Lock()
	defer version_manager.mutex.Unlock()
	return version_manager.root
}

func (version_manager *VersionManager) NewChild(parent *Version) *Version {
	version_manager.mutex.Lock()
	defer version_manager.mutex.Unlock()

	if parent == nil {
		parent = version_manager.root
	}

	version := &Version{
		id:     version_manager.next,
		parent: parent,
	}
	version_manager.next++
	return version
}

func IsAncestor(candidate, version *Version) bool {
	if candidate == nil || version == nil {
		return false
	}
	for current := version; current != nil; current = current.parent {
		if current == candidate {
			return true
		}
	}
	return false
}

func CommonAncestor(a, b *Version) *Version {
	if a == nil || b == nil {
		return nil
	}

	depth_a := depth(a)
	depth_b := depth(b)

	for depth_a > depth_b {
		a = a.parent
		depth_a--
	}
	for depth_b > depth_a {
		b = b.parent
		depth_b--
	}

	for a != nil && b != nil {
		if a == b {
			return a
		}
		a = a.parent
		b = b.parent
	}
	return nil
}

func depth(v *Version) int {
	d := 0
	for cur := v; cur != nil; cur = cur.parent {
		d++
	}
	return d - 1
}

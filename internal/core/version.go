package core

import "sync"

type VersionID uint64

type Version struct {
	id     VersionID
	parent *Version
}

func (version *Version) ID() VersionID {
	return version.id
}

func (version *Version) Parent() *Version {
	return version.parent
}

type VersionManager struct {
	mutex sync.Mutex
	next  VersionID
	root  *Version
}

func NewVersionManager() *VersionManager {
	versionManager := &VersionManager{next: 1}
	versionManager.root = &Version{id: versionManager.next}
	versionManager.next++
	return versionManager
}

func (v *VersionManager) Root() *Version {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	return v.root
}

func (v *VersionManager) NewChild(parent *Version) *Version {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	if parent == nil {
		parent = v.root
	}

	version := &Version{
		id:     v.next,
		parent: parent,
	}
	v.next++
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

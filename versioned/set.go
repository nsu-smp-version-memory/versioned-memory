package versioned

import "github.com/nsu-smp-version-memory/versioned-memory/internal/core"

type TagSet map[core.OperationID]struct{}

func (set TagSet) Copy() TagSet {
	out := make(TagSet, len(set))
	for k := range set {
		out[k] = struct{}{}
	}
	return out
}

type node struct {
	value int

	adds TagSet
	dels TagSet

	left  *node
	right *node
}

type Set struct {
	root    *node
	version *core.Version
}

func New(version *core.Version) Set {
	return Set{version: version}
}

func (set Set) Version() *core.Version {
	return set.version
}

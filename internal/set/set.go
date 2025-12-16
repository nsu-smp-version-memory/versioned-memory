package set

import "core"

type TagSet map[core.OperationID]struct{}

func (s TagSet) Clone() TagSet {
	if len(s) == 0 {
		return nil
	}
	out := make(TagSet, len(s))
	for k := range s {
		out[k] = struct{}{}
	}
	return out
}

type node struct {
	key int

	adds TagSet
	dels TagSet

	left  *node
	right *node
}

type Set struct {
	root *node
	v    *core.Version
}

func New(v *core.Version) Set {
	return Set{v: v}
}

func (s Set) Version() *core.Version {
	return s.v
}

package timeline

import (
	"github.com/nsu-smp-version-memory/versioned-memory/internal/version"
)

type Node[DIFF any] struct {
	version uint64
	prev    *Node[DIFF]
	diff    DIFF
}

type Timeline[DIFF any] struct {
	last          *Node[DIFF]
	versionSource *version.Source
}

func New[DIFF any]() *Timeline[DIFF] {
	t := &Timeline[DIFF]{
		last:          nil,
		versionSource: version.NewSource(),
	}

	return t
}

func (t *Timeline[DIFF]) NextChange(diff DIFF) *Timeline[DIFF] {
	last := &Node[DIFF]{
		version: t.versionSource.Peek(),
		prev:    t.last,
		diff:    diff,
	}

	return &Timeline[DIFF]{
		last:          last,
		versionSource: t.versionSource,
	}
}

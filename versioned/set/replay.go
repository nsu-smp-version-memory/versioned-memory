package set

import (
	"sort"

	"github.com/nsu-smp-version-memory/versioned-memory/internal/core"
)

func replayToMap(tl *core.Timeline[Diff]) map[int]struct{} {
	out := make(map[int]struct{})

	for _, ops := range tl.Operations() {
		switch ops.Diff.Kind {
		case Add:
			out[ops.Diff.Value] = struct{}{}
		case Remove:
			delete(out, ops.Diff.Value)
		}
	}
	return out
}

func mapKeysSorted(m map[int]struct{}) []int {
	values := make([]int, 0, len(m))
	for v := range m {
		values = append(values, v)
	}
	sort.Ints(values)
	return values
}

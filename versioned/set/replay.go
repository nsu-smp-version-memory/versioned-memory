package set

import (
	"sort"

	"github.com/nsu-smp-version-memory/versioned-memory/internal/core"
)

func replayToMap(tl *core.Timeline[Diff]) map[int]struct{} {
	present := make(map[int]struct{})

	operations := tl.Operations()

	sort.Slice(operations, func(i, j int) bool { return operations[i].ID.Before(operations[j].ID) })

	for _, ops := range operations {
		switch ops.Diff.Kind {
		case Add:
			present[ops.Diff.Value] = struct{}{}
		case Remove:
			delete(present, ops.Diff.Value)
		}
	}
	return present
}

func mapKeysSorted(m map[int]struct{}) []int {
	values := make([]int, 0, len(m))
	for v := range m {
		values = append(values, v)
	}
	sort.Ints(values)
	return values
}

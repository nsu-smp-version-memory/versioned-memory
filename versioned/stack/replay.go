package stack

import "github.com/nsu-smp-version-memory/versioned-memory/internal/core"

func replayToSlice(tl *core.Timeline[Diff]) []int {

	ops := tl.Operations()

	out := make([]int, 0, len(ops))
	for _, op := range ops {
		switch op.Diff.Kind {
		case Push:
			out = append(out, op.Diff.Value)
		case Pop:
			if len(out) > 0 {
				out = out[:len(out)-1]
			}
		}
	}
	return out
}

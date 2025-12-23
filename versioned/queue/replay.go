package queue

import "github.com/nsu-smp-version-memory/versioned-memory/internal/core"

func replayToSlice(tl *core.Timeline[Diff]) []int {
	ops := tl.Operations()

	out := make([]int, 0, len(ops))
	head := 0

	for _, op := range ops {
		switch op.Diff.Kind {
		case Enqueue:
			out = append(out, op.Diff.Value)
		case Dequeue:
			if head < len(out) {
				head++
			}
		}
	}

	if head == 0 {
		return out
	}
	if head >= len(out) {
		return nil
	}
	alive := make([]int, len(out)-head)
	copy(alive, out[head:])
	return alive
}

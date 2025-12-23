package tree

import (
	"github.com/nsu-smp-version-memory/versioned-memory/internal/core"
)

type pendingBranch struct {
	timeline *core.Timeline[Diff]
	since    core.ForkPoint[Diff]
}

func (t *Tree) WithBranch(fn func(local *Tree)) <-chan struct{} {
	done := make(chan struct{})

	t.mutex.Lock()
	branchTimeline, since := t.timeline.Fork(core.NewSource())
	merger := t.merger
	t.mutex.Unlock()

	local := &Tree{
		timeline: branchTimeline,
		merger:   merger,
	}

	t.wg.Go(func() {
		defer close(done)

		fn(local)

		t.mutex.Lock()
		t.pendingBranches = append(t.pendingBranches, pendingBranch{
			timeline: local.timeline,
			since:    since,
		})
		t.mutex.Unlock()
	})

	return done
}

func (t *Tree) MergeBranches() {
	t.mutex.Lock()
	pending := t.pendingBranches
	t.pendingBranches = nil
	base := t.timeline
	root := t.root
	t.mutex.Unlock()

	input := make([][]core.Operation[Diff], 0)

	if base != nil {
		operations := base.Operations()
		for i := len(operations) - 1; i >= 0; i-- {
			switch operations[i].Diff.Kind {
			case Add:
				root = root.remove(operations[i].Diff.Value)
			case Remove:
				root = root.insert(operations[i].Diff.Value)
			}

		}

		input = append(input, operations)
	}

	for _, br := range pending {
		if br.timeline == nil {
			continue
		}
		input = append(input, br.timeline.OperationsAfter(br.since))
	}

	result := t.merger.Merge(input)

	for _, op := range result {
		switch op.Diff.Kind {
		case Add:
			root = root.insert(op.Diff.Value)
		case Remove:
			root = root.remove(op.Diff.Value)
		}
	}

	t.mutex.Lock()
	t.timeline = core.TimelineFromOperations(core.NewSource(), result)
	t.root = root
	t.mutex.Unlock()
}

func (t *Tree) Go(f func(t *Tree)) {
	t.WithBranch(f)
}

func (t *Tree) Join() {
	t.wg.Wait()
	t.MergeBranches()
}

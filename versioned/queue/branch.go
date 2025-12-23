package queue

import "github.com/nsu-smp-version-memory/versioned-memory/internal/core"

type pendingBranch struct {
	timeline *core.Timeline[Diff]
	since    core.ForkPoint[Diff]
}

func (q *Queue) WithBranch(fn func(local *Queue)) <-chan struct{} {
	done := make(chan struct{})

	q.mutex.Lock()
	branchTimeline, since := q.timeline.Fork(core.NewSource())

	merger := q.merger
	q.mutex.Unlock()

	local := &Queue{
		timeline: branchTimeline,
		merger:   merger,
	}

	q.wg.Go(func() {
		defer close(done)

		fn(local)

		q.mutex.Lock()
		q.pendingBranches = append(q.pendingBranches, pendingBranch{
			timeline: local.timeline,
			since:    since,
		})
		q.mutex.Unlock()
	})

	return done
}

func (q *Queue) MergeBranches() {
	q.mutex.Lock()
	pending := q.pendingBranches
	q.pendingBranches = nil
	base := q.timeline
	merger := q.merger
	q.mutex.Unlock()

	input := make([][]core.Operation[Diff], 0)

	if base != nil {
		input = append(input, base.Operations())
	}

	for _, br := range pending {
		input = append(input, br.timeline.OperationsAfter(br.since))
	}

	result := merger.Merge(input)

	q.mutex.Lock()
	q.timeline = core.TimelineFromOperations(core.NewSource(), result)
	q.mutex.Unlock()
}

func (q *Queue) Go(f func(q *Queue)) {
	q.WithBranch(f)
}

func (q *Queue) Join() {
	q.wg.Wait()
	q.MergeBranches()
}

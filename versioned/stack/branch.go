package stack

import "github.com/nsu-smp-version-memory/versioned-memory/internal/core"

type pendingBranch struct {
	timeline *core.Timeline[Diff]
	since    core.ForkPoint[Diff]
}

func (s *Stack) WithBranch(fn func(local *Stack)) <-chan struct{} {
	done := make(chan struct{})

	s.mutex.Lock()
	branchTimeline, since := s.timeline.Fork(core.NewSource())
	merger := s.merger
	s.mutex.Unlock()

	local := &Stack{
		timeline: branchTimeline,
		merger:   merger,
	}

	s.wg.Go(func() {
		defer close(done)

		fn(local)

		s.mutex.Lock()
		s.pendingBranches = append(s.pendingBranches, pendingBranch{
			timeline: local.timeline,
			since:    since,
		})
		s.mutex.Unlock()
	})

	return done
}

func (s *Stack) MergeBranches() {
	s.mutex.Lock()
	pending := s.pendingBranches
	s.pendingBranches = nil
	base := s.timeline
	merger := s.merger
	s.mutex.Unlock()

	input := make([][]core.Operation[Diff], 0)

	if base != nil {
		input = append(input, base.Operations())
	}

	for _, br := range pending {
		if br.timeline == nil {
			continue
		}
		input = append(input, br.timeline.OperationsAfter(br.since))
	}

	result := merger.Merge(input)

	s.mutex.Lock()
	s.timeline = core.TimelineFromOperations(core.NewSource(), result)
	s.mutex.Unlock()
}

func (s *Stack) Go(f func(s *Stack)) {
	s.WithBranch(f)
}

func (s *Stack) Join() {
	s.wg.Wait()
	s.MergeBranches()
}

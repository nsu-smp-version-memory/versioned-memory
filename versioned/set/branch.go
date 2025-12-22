package set

import (
	"github.com/nsu-smp-version-memory/versioned-memory/internal/core"
)

type pendingBranch struct {
	timeline *core.Timeline[Diff]
	since    core.ForkPoint[Diff]
}

func (s *Set) WithBranch(fn func(local *Set)) <-chan struct{} {
	done := make(chan struct{})

	s.mutex.Lock()
	branchTimeline, since := s.timeline.Fork(core.NewSource())
	s.mutex.Unlock()

	local := &Set{timeline: branchTimeline}

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

func (s *Set) MergeBranches() {
	s.mutex.Lock()
	pending := s.pendingBranches
	s.pendingBranches = nil
	base := s.timeline
	s.mutex.Unlock()

	ops := make([]core.Operation[Diff], 0)

	if base != nil {
		ops = append(ops, base.Operations()...)
	}

	for _, br := range pending {
		if br.timeline == nil {
			continue
		}
		ops = append(ops, br.timeline.OperationsAfter(br.since)...)
	}

	sortOperationsByID(ops)

	s.mutex.Lock()
	s.timeline = core.TimelineFromOperations(core.NewSource(), ops)
	s.mutex.Unlock()
}

func (s *Set) Go(f func(s *Set)) {
	s.WithBranch(f)
}

func (s *Set) Join() {
	s.wg.Wait()
	s.MergeBranches()
}

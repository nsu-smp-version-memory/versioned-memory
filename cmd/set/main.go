package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/nsu-smp-version-memory/versioned-memory/internal/core"
	"github.com/nsu-smp-version-memory/versioned-memory/pkg/versioned"
)

func main() {
	versionManager := core.NewVersionManager()
	rootV := versionManager.Root()

	workers := 4
	opsPerWorker := 30
	keySpace := 20

	type result struct {
		id  int
		set set.Set
	}

	results := make([]result, workers)
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		i := i
		go func() {
			defer wg.Done()

			src := core.NewSource(core.SourceID(i+1), versionManager)
			local := set.New(rootV)

			r := rand.New(rand.NewSource(int64(1000 + i)))

			for j := 0; j < opsPerWorker; j++ {
				time.Sleep(time.Duration(r.Intn(35)) * time.Millisecond)

				k := r.Intn(keySpace)
				if r.Intn(100) < 70 {
					local = local.Add(src, k)
				} else {
					local = local.Remove(src, k)
				}
			}

			results[i] = result{id: i + 1, set: local}
		}()
	}

	wg.Wait()

	fmt.Println("=== per-worker states ===")
	for _, res := range results {
		fmt.Printf("worker %d: version=%d keys=%v height=%d\n",
			res.id,
			res.set.Version().ID(),
			res.set.Keys(),
			res.set.Height(),
		)
	}

	merged := results[0].set
	for i := 1; i < workers; i++ {
		merged = merged.Merge(versionManager, results[i].set)
	}

	fmt.Println("\n=== merged state ===")
	fmt.Printf("merged version=%d keys=%v height=%d\n",
		merged.Version().ID(),
		merged.Keys(),
		merged.Height(),
	)
}

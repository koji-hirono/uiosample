package worker

import (
	"context"
	"runtime"
	"sync"

	"golang.org/x/sys/unix"
)

var lock bool = true

func Worker(ctx context.Context, wg *sync.WaitGroup, cpu int, task func()) {
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		if lock {
			runtime.LockOSThread()
			var cpuset unix.CPUSet
			cpuset.Set(cpu)
			err := unix.SchedSetaffinity(0, &cpuset)
			if err != nil {
				return
			}
		}
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			task()
		}
	}(ctx)
}

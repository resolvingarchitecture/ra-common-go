package tasks

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestRunnerOneShotAndPeriodic(t *testing.T) {
	runner := NewRunner(20)
	var runs atomic.Int32
	runner.AddTask(NewLambdaTask(Once("c"), func() bool { runs.Add(1); return true }))
	runner.Start()
	for i := 0; i < 100 && (runs.Load() < 1 || runner.TaskCount() > 0); i++ {
		time.Sleep(10 * time.Millisecond)
	}
	if runs.Load() != 1 {
		t.Errorf("runs = %d, want 1", runs.Load())
	}
	runner.Shutdown()
	if runner.Status() != RunnerShutdown {
		t.Errorf("Status() = %q, want Shutdown", runner.Status())
	}

	runner2 := NewRunner(10)
	var n atomic.Int32
	runner2.AddTask(NewLambdaTask(Periodic("p", 10), func() bool { n.Add(1); return true }))
	runner2.Start()
	time.Sleep(150 * time.Millisecond)
	runner2.Shutdown()
	if n.Load() < 2 {
		t.Errorf("n = %d, want >= 2", n.Load())
	}
}

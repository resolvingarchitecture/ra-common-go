package tasks

import (
	"sync"
	"time"
)

type RunnerStatus string

const (
	RunnerRunning  RunnerStatus = "Running"
	RunnerStopping RunnerStatus = "Stopping"
	RunnerShutdown RunnerStatus = "Shutdown"
)

type managedTask struct {
	task      Task
	status    Status
	scheduled bool
}

// Runner schedules and runs Task objects. Not safe to Start() again after
// Shutdown() returns - create a new Runner instead.
type Runner struct {
	mu      sync.Mutex
	tasks   []*managedTask
	running bool
	pollMs  time.Duration
	stopCh  chan struct{}
	wg      sync.WaitGroup
}

func NewRunner(pollPeriodMs int) *Runner {
	if pollPeriodMs < 1 {
		pollPeriodMs = 1
	}
	return &Runner{pollMs: time.Duration(pollPeriodMs) * time.Millisecond}
}

func (r *Runner) AddTask(t Task) {
	r.mu.Lock()
	r.tasks = append(r.tasks, &managedTask{task: t, status: Ready})
	running := r.running
	r.mu.Unlock()
	if running {
		r.pass()
	}
}

func (r *Runner) Poke() {
	r.mu.Lock()
	running := r.running
	r.mu.Unlock()
	if running {
		r.pass()
	}
}

func (r *Runner) Start() {
	r.mu.Lock()
	if r.running {
		r.mu.Unlock()
		return
	}
	r.running = true
	r.stopCh = make(chan struct{})
	r.mu.Unlock()

	r.wg.Add(1)
	go r.pollLoop()
	r.pass()
}

func (r *Runner) pollLoop() {
	defer r.wg.Done()
	ticker := time.NewTicker(r.pollMs)
	defer ticker.Stop()
	for {
		select {
		case <-r.stopCh:
			return
		case <-ticker.C:
			r.pass()
		}
	}
}

func (r *Runner) Status() RunnerStatus {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.running && r.stopCh == nil {
		return RunnerShutdown
	}
	if r.running {
		return RunnerRunning
	}
	return RunnerStopping
}

func (r *Runner) TaskCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.tasks)
}

// Shutdown stops the poll loop and waits for every in-flight task worker to finish.
func (r *Runner) Shutdown() {
	r.mu.Lock()
	if r.running {
		r.running = false
		close(r.stopCh)
	}
	r.mu.Unlock()

	r.wg.Wait()

	r.mu.Lock()
	r.tasks = nil
	r.stopCh = nil
	r.mu.Unlock()
}

func (r *Runner) pass() {
	r.mu.Lock()
	filtered := r.tasks[:0]
	for _, m := range r.tasks {
		if m.status != Completed {
			filtered = append(filtered, m)
		}
	}
	r.tasks = filtered

	var toSchedule []*managedTask
	for _, m := range r.tasks {
		if m.scheduled {
			continue
		}
		cfg := m.task.Config()
		if cfg.PeriodicityMs == -1 {
			continue
		}
		m.scheduled = true
		toSchedule = append(toSchedule, m)
	}
	stopCh := r.stopCh
	r.mu.Unlock()

	for _, m := range toSchedule {
		r.wg.Add(1)
		go r.worker(m, stopCh)
	}
}

func (r *Runner) worker(m *managedTask, stopCh chan struct{}) {
	defer r.wg.Done()
	cfg := m.task.Config()
	if cfg.Delayed && cfg.DelayMs > 0 {
		select {
		case <-time.After(time.Duration(cfg.DelayMs) * time.Millisecond):
		case <-stopCh:
		}
	}
	for {
		select {
		case <-stopCh:
			m.task.OnStop()
			r.setStatus(m, Completed)
			return
		default:
		}
		if m.task.ShouldStop() {
			m.task.OnStop()
			r.setStatus(m, Completed)
			return
		}
		r.setStatus(m, Running)
		m.task.Execute()
		if cfg.PeriodicityMs <= 0 {
			r.setStatus(m, Completed)
			return
		}
		r.setStatus(m, Ready)
		select {
		case <-time.After(time.Duration(cfg.PeriodicityMs) * time.Millisecond):
		case <-stopCh:
			r.setStatus(m, Completed)
			return
		}
	}
}

func (r *Runner) setStatus(m *managedTask, s Status) {
	r.mu.Lock()
	m.status = s
	r.mu.Unlock()
}

// Package tasks: background task scheduling. Ports the ra.common.tasks
// package. Built on goroutines + channels rather than the raw std::thread
// juggling ra-common-cpp needed - Go's scheduler and channel-based shutdown
// signaling make the concurrency here much lower-risk (see that port's
// DESIGN.md for what goes wrong without them).
package tasks

type Status string

const (
	Ready     Status = "Ready"
	Running   Status = "Running"
	Completed Status = "Completed"
)

type Config struct {
	Name          string
	PeriodicityMs int
	Delayed       bool
	DelayMs       int
	FixedDelay    bool
	LongRunning   bool
}

func Once(name string) Config { return Config{Name: name} }

func Periodic(name string, periodMs int) Config {
	c := Once(name)
	c.PeriodicityMs = periodMs
	return c
}

// Task is a schedulable unit of work.
type Task interface {
	Config() Config
	Execute() bool
	ShouldStop() bool
	OnStop()
}

// LambdaTask wraps a config + a func() bool as a Task, for the common case
// of not wanting a full type.
type LambdaTask struct {
	config  Config
	execute func() bool
}

func NewLambdaTask(config Config, execute func() bool) *LambdaTask {
	return &LambdaTask{config: config, execute: execute}
}

func (t *LambdaTask) Config() Config   { return t.config }
func (t *LambdaTask) Execute() bool    { return t.execute() }
func (t *LambdaTask) ShouldStop() bool { return false }
func (t *LambdaTask) OnStop()          {}

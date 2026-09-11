package service

import (
	"testing"

	"github.com/resolvingarchitecture/ra-common-go/messaging"
	"github.com/resolvingarchitecture/ra-common-go/servicestatus"
)

type toy struct {
	BaseService
	core    *ServiceCore
	started bool
}

func newToy() *toy { return &toy{core: NewServiceCore("ra.test.Toy")} }

func (t *toy) Core() *ServiceCore { return t.core }
func (t *toy) Start(map[string]string) bool {
	t.started = true
	t.core.UpdateStatus(servicestatus.Running)
	return true
}
func (t *toy) Shutdown() bool                      { t.started = false; return true }
func (t *toy) GracefulShutdown() bool              { return DefaultGracefulShutdown(t) }
func (t *toy) HandleCommand(e *messaging.Envelope) { DefaultHandleCommand(t, e) }

func TestServiceCommandDrivesLifecycle(t *testing.T) {
	svc := newToy()
	start := messaging.CmdStart
	e := messaging.CommandEnvelope()
	e.MessageValue = messaging.NewCommandMessage(&start)
	Handle(svc, e)
	if !svc.started {
		t.Error("Start should have run")
	}
	if ServiceStatusValue(svc) != servicestatus.Running {
		t.Errorf("status = %q, want Running", ServiceStatusValue(svc))
	}

	report := messaging.CmdReport
	r := messaging.CommandEnvelope()
	r.MessageValue = messaging.NewCommandMessage(&report)
	Handle(svc, r)
	if !r.HeaderExists("result") {
		t.Error("Report command should set the result header")
	}
}

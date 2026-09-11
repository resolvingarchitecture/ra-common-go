// Package service: the ServiceCore shared state and Service interface. Ports
// ra.common.service.{Service, BaseService} (except ServiceDaemon, deferred).
package service

import (
	"github.com/resolvingarchitecture/ra-common-go/messaging"
	"github.com/resolvingarchitecture/ra-common-go/servicestatus"
)

// ServiceCore holds state shared by every service (ports the BaseService fields).
type ServiceCore struct {
	ServiceClassName      string
	StatusValue           servicestatus.ServiceStatus
	Registered            bool
	Version               *string
	ServicesDependentUpon []string
	Config                map[string]string
	Producer              messaging.MessageProducer // set by the hosting bus; nil until then
	Observer              servicestatus.ServiceStatusObserver
}

func NewServiceCore(serviceClassName string) *ServiceCore {
	return &ServiceCore{
		ServiceClassName: serviceClassName,
		StatusValue:      servicestatus.NotInitialized,
		Config:           map[string]string{},
	}
}

func (c *ServiceCore) AddDependentService(name string) {
	c.ServicesDependentUpon = append(c.ServicesDependentUpon, name)
}

func (c *ServiceCore) Send(envelope *messaging.Envelope) bool {
	if c.Producer == nil {
		return false
	}
	return c.Producer.Send(envelope)
}

func (c *ServiceCore) Report() servicestatus.ServiceReport {
	return servicestatus.ServiceReport{
		ServiceClassName:      c.ServiceClassName,
		ServiceStatusValue:    c.StatusValue,
		Registered:            c.Registered,
		Running:               c.StatusValue == servicestatus.Running,
		Version:               c.Version,
		ServicesDependentUpon: c.ServicesDependentUpon,
	}
}

func (c *ServiceCore) UpdateStatus(status servicestatus.ServiceStatus) {
	if c.StatusValue == status {
		return
	}
	c.StatusValue = status
	if c.Observer != nil {
		c.Observer.ServiceStatusChanged(c.ServiceClassName, status)
	}
	if c.Producer != nil {
		ev := messaging.NewEventMessage(messaging.EventTypeServiceStatus)
		ev.MessageValue = c.Report()
		e := messaging.EventEnvelope(messaging.EventTypeServiceStatus)
		e.MessageValue = ev
		e.AddRoute("ra.notification.NotificationService", "PUBLISH")
		e.Ratchet()
		c.Send(e)
	}
}

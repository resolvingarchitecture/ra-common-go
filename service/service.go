package service

import (
	"github.com/resolvingarchitecture/ra-common-go/lifecycle"
	"github.com/resolvingarchitecture/ra-common-go/messaging"
	"github.com/resolvingarchitecture/ra-common-go/servicestatus"
)

// Service is a message-driven service. Ports ra.common.service.Service.
//
// Go has no virtual dispatch through struct embedding - a method BaseService
// provides can't call back into a concrete type's override the way an
// abstract class's template method can in every other port. So
// GracefulShutdown/HandleCommand/the top-level Handle dispatcher are free
// functions taking the full Service interface (real dynamic dispatch through
// an interface value) rather than BaseService methods; a concrete service
// calls the matching Default* function from its own one-line method body.
// See DESIGN.md.
type Service interface {
	lifecycle.LifeCycle
	Core() *ServiceCore
	HandleDocument(envelope *messaging.Envelope)
	HandleEvent(envelope *messaging.Envelope)
	HandleCommand(envelope *messaging.Envelope)
	HandleHeaders(envelope *messaging.Envelope)
}

// BaseService supplies the parts of Service that never need to call back
// into an overridden method. Embed it and implement Start/Shutdown/Core
// yourself, plus one-line HandleCommand/GracefulShutdown methods delegating
// to DefaultHandleCommand/DefaultGracefulShutdown below.
type BaseService struct{}

func (BaseService) Pause() bool                        { return false }
func (BaseService) Unpause() bool                      { return false }
func (BaseService) Restart() bool                      { return false }
func (BaseService) HandleDocument(*messaging.Envelope) {}
func (BaseService) HandleEvent(*messaging.Envelope)    {}
func (BaseService) HandleHeaders(*messaging.Envelope)  {}

// DefaultGracefulShutdown is ra.common.service.Service's default
// GracefulShutdown(); call it from a concrete service's own method.
func DefaultGracefulShutdown(s Service) bool { return s.Shutdown() }

// DefaultHandleCommand is ra.common.service.Service's default
// HandleCommand(envelope); call it from a concrete service's own method.
func DefaultHandleCommand(s Service, envelope *messaging.Envelope) {
	cmdMsg, ok := envelope.MessageValue.(*messaging.CommandMessage)
	if !ok || cmdMsg.CommandValue == nil {
		return
	}
	config := make(map[string]string, len(s.Core().Config))
	for k, v := range s.Core().Config {
		config[k] = v
	}
	switch *cmdMsg.CommandValue {
	case messaging.CmdStart:
		s.Start(config)
	case messaging.CmdPause:
		s.Pause()
	case messaging.CmdUnpause:
		s.Unpause()
	case messaging.CmdRestart:
		s.Restart()
	case messaging.CmdShutdown:
		s.Shutdown()
	case messaging.CmdGracefullyShutdown:
		s.GracefulShutdown()
	case messaging.CmdReport:
		envelope.SetHeader("result", s.Core().Report())
	}
}

func ServiceStatusValue(s Service) servicestatus.ServiceStatus { return s.Core().StatusValue }

func Report(s Service) servicestatus.ServiceReport { return s.Core().Report() }

// Handle dispatches envelope to the matching Handle* method based on its message kind.
func Handle(s Service, envelope *messaging.Envelope) *messaging.Envelope {
	switch envelope.MessageValue.(type) {
	case *messaging.DocumentMessage:
		s.HandleDocument(envelope)
	case *messaging.EventMessage:
		s.HandleEvent(envelope)
	case *messaging.CommandMessage:
		s.HandleCommand(envelope)
	default:
		s.HandleHeaders(envelope)
	}
	return envelope
}

// Package servicestatus holds service status/level enums and the report
// shape. A leaf package so messaging.Envelope can depend on ServiceLevel
// without pulling in the full service package (which depends on messaging).
// Ports ra.common.service.{ServiceStatus, ServiceLevel, ServiceReport,
// ServiceMessage, ServiceStatusObserver}.
package servicestatus

type ServiceLevel string

const (
	AtMostOnce  ServiceLevel = "AtMostOnce"
	AtLeastOnce ServiceLevel = "AtLeastOnce"
	ExactlyOnce ServiceLevel = "ExactlyOnce"
)

type ServiceStatus string

const (
	NotInitialized         ServiceStatus = "NotInitialized"
	Initializing           ServiceStatus = "Initializing"
	Waiting                ServiceStatus = "Waiting"
	Starting               ServiceStatus = "Starting"
	Running                ServiceStatus = "Running"
	Verified               ServiceStatus = "Verified"
	PartiallyRunning       ServiceStatus = "PartiallyRunning"
	DegradedRunning        ServiceStatus = "DegradedRunning"
	Unstable               ServiceStatus = "Unstable"
	Pausing                ServiceStatus = "Pausing"
	Paused                 ServiceStatus = "Paused"
	Unpausing              ServiceStatus = "Unpausing"
	ShuttingDown           ServiceStatus = "ShuttingDown"
	GracefullyShuttingDown ServiceStatus = "GracefullyShuttingDown"
	Shutdown               ServiceStatus = "Shutdown"
	GracefullyShutdown     ServiceStatus = "GracefullyShutdown"
	Restarting             ServiceStatus = "Restarting"
	Unavailable            ServiceStatus = "Unavailable"
	Error                  ServiceStatus = "Error"
)

var runningStates = map[ServiceStatus]bool{
	Running: true, Verified: true, PartiallyRunning: true, DegradedRunning: true,
}

func IsRunning(status ServiceStatus) bool { return runningStates[status] }

const (
	NoError         = -1
	RequestRequired = 0
	RaServiceImpl   = "ra.service.impl"
)

type ServiceMessage struct {
	StatusCode   int     `json:"status_code"`
	ErrorMessage *string `json:"error_message,omitempty"`
	Exception    *string `json:"exception,omitempty"`
	Type         *string `json:"type,omitempty"`
}

type ServiceReport struct {
	ServiceClassName      string        `json:"service_class_name"`
	ServiceStatusValue    ServiceStatus `json:"service_status"`
	Registered            bool          `json:"registered"`
	Running               bool          `json:"running"`
	Version               *string       `json:"version,omitempty"`
	ServicesDependentUpon []string      `json:"services_dependent_upon,omitempty"`
}

type ServiceStatusObserver interface {
	ServiceStatusChanged(serviceFullName string, status ServiceStatus)
}

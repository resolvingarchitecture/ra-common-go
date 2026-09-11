// Package lifecycle holds the Status enum and LifeCycle interface. Ports
// ra.common.Status and ra.common.LifeCycle. A leaf package (no Envelope
// dependency) so messaging's channel/bus interfaces can embed LifeCycle
// without a cycle - ra.common.Client (which does need Envelope) lives in
// messaging instead, next to Envelope itself.
package lifecycle

type Status string

const (
	Initialized Status = "Initialized"
	Starting    Status = "Starting"
	Running     Status = "Running"
	Paused      Status = "Paused"
	Stopping    Status = "Stopping"
	Stopped     Status = "Stopped"
	Errored     Status = "Errored"
)

// LifeCycle is the start/pause/restart/shutdown contract. Every method
// returns true on success, matching the Java API. Unpause is named as such
// (rather than Resume) to mirror the Java note about the Thread.resume clash.
type LifeCycle interface {
	Start(properties map[string]string) bool
	Pause() bool
	Unpause() bool
	Restart() bool
	Shutdown() bool
	GracefulShutdown() bool
}

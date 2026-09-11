package messaging

import "encoding/json"

type Command string

const (
	CmdStart                         Command = "Start"
	CmdShutdown                      Command = "Shutdown"
	CmdGracefullyShutdown            Command = "GracefullyShutdown"
	CmdRestart                       Command = "Restart"
	CmdPause                         Command = "Pause"
	CmdUnpause                       Command = "Unpause"
	CmdNetState                      Command = "NetState"
	CmdReport                        Command = "Report"
	CmdRegisterStateChangeListener   Command = "RegisterStateChangeListener"
	CmdUnregisterStateChangeListener Command = "UnregisterStateChangeListener"
)

type CommandMessage struct {
	BaseMessage
	CommandValue *Command `json:"command,omitempty"`
}

func NewCommandMessage(cmd *Command) *CommandMessage { return &CommandMessage{CommandValue: cmd} }

func (m *CommandMessage) Kind() string { return "command" }

func (m CommandMessage) MarshalJSON() ([]byte, error) {
	type alias CommandMessage
	return json.Marshal(struct {
		Kind string `json:"kind"`
		alias
	}{Kind: "command", alias: alias(m)})
}

package lights

import (
	"fmt"
	"strings"
)

// Command represents an Inception system command. Parts are split using `|`
// characters.
type Command struct {
	Action string
	Type   string
	ID     string
	Parts  []string
}

// NewCommand parses a command string to obtain it's command, message, and ID.
// An error is returned if the command or message is not recognized.
func NewCommand(cmd string) (*Command, error) {
	if len(cmd) < 2 {
		// Command is too small
		return nil, fmt.Errorf("Truncated command received: %s", cmd)
	}
	command := &Command{}
	switch cmd[0] {
	case '!': // Execute command
		command.Action = "execute"
	case '+': // Add configuration
		command.Action = "add"
	case '-': // Remove configuration
		command.Action = "remove"
	case '?':
		command.Action = "query"
	default:
		return nil, fmt.Errorf("Unknown action code '%s' in command: %s", string(cmd[0]), cmd)
	}
	switch cmd[1] {
	case '#': // Color
		command.Type = "color"
	case ':': // Pattern
		command.Type = "pattern"
	case '~': // Schedule
		command.Type = "schedule"
	case '^': // Scene
		command.Type = "scene"
	case '-': // Property
		command.Type = "property"
	default:
		return nil, fmt.Errorf("Unknown type code '%s' in command: %s", string(cmd[1]), cmd)
	}
	if len(cmd) > 2 {
		command.Parts = strings.Split(cmd[2:], "|")
		switch command.Type {
		case "color":
			command.ID = "#" + command.Parts[0]
		case "pattern":
			command.ID = strings.Split(command.Parts[0], ":")[0]
		case "schedule":
			command.ID = command.Parts[0]
		case "scene":
			command.ID = command.Parts[0]
		case "property":
			command.ID = command.Parts[0]
		}
	} else {
		command.Parts = []string{}
	}

	return command, nil
}

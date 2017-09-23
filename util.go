package lights

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

// PanicIf checks if the provided error is non-nil and panics.
// If messages are provided they are logged whether the error occurred or not.
func PanicIf(err error, messages ...string) {
	if err != nil {
		if len(messages) > 0 {
			log.Println("PANIC", strings.Join(messages, " "), err)
		} else {
			log.Println("PANIC", err)
		}
		panic(err)
	}
	if len(messages) > 0 {
		log.Println(strings.Join(messages, " "))
	}
}

// CheckIf returns true if the provided error is non-nil.
// If messages are provided they are logged whether the error occurred or not.
func CheckIf(err error, messages ...string) bool {
	if err != nil {
		if len(messages) > 0 {
			log.Println("ERROR", strings.Join(messages, " "), err)
		} else {
			log.Println("ERROR", err)
		}
		return true
	}
	if len(messages) > 0 {
		log.Println(strings.Join(messages, " "))
	}
	return false
}

// ReadFileOrPanic reads in a file or panics if there was a problem.
func ReadFileOrPanic(path string) []byte {
	data, err := ioutil.ReadFile(path)
	PanicIf(err)
	return data
}

// ReadFileAsStringOrPanic reads a file content as utf-8 encoded string
// or panics if there was a problem.
func ReadFileAsStringOrPanic(path string) string {
	return string(ReadFileOrPanic(path))
}

// PrepPath prepares a path for use by cleaning and converting to an absolute
// path.
func PrepPath(path string) (string, error) {
	// Convert to absolute path - should make file-not-found errors report
	// the full path of the file that was attempted to be read.
	return filepath.Abs(filepath.Clean(path))
}

// AgentPort looks up the correct port for an agent by name.
func AgentPort(agent string) string {
	switch agent {
	case "gateway":
		return "8000"
	case "controller":
		return "8002"
	case "gatekeeper":
		return "8003"
	case "scheduler":
		return "8004"
	case "updater":
		return "8005"
	default:
		// Default is the gateway agent
		return ""
	}
}

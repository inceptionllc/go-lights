package lights

import (
	"io/ioutil"
	"log"
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

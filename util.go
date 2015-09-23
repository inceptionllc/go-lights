package lights

import (
	"io/ioutil"
	"log"
)

// PanicIf checks if the provided error is non-nil and panics.
func PanicIf(err error, messages ...string) {
	if err != nil {
		if len(messages) > 0 {
			log.Println("PANIC", messages)
		}
		panic(err)
	}
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

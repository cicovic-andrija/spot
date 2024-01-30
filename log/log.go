package log

import (
	"log"
	"os"
)

var (
	info  = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
)

// Fatalf calls Fatalf of the underlying log
func Fatalf(format string, a ...interface{}) {
	error.Fatalf(format, a...)
}

// Errorf calls Printf of the underlying log
func Errorf(format string, a ...interface{}) {
	error.Printf(format, a...)
}

// Infof calls Printf of the underlying log
func Infof(format string, a ...interface{}) {
	info.Printf(format, a...)
}

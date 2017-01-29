package encoding

import (
	"fmt"
	"os"
)

// Package level logger initialized with the s
var logger = Logger(log{})

// Logger is an interface for a logger
type Logger interface {
	Printf(string, ...interface{})
}

// Simple logger
type log struct{}

func (l log) Printf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

// SetLogger set the logger for the package
func SetLogger(l Logger) {
	logger = l
}

package logging

import (
	"fmt"

	"github.com/fatih/color"
)

// Logger is a struct that represents a logger.
type Logger struct {
	debug bool
}

// NewLogger creates a new Logger instance.
func NewLogger(debug bool) *Logger {
	return &Logger{
		debug: debug,
	}
}

// Log logs a message to the console.
func (l *Logger) Log(message string) {
	fmt.Println(color.WhiteString(message))
}

// Debug logs a message to the console if debug is enabled.
func (l *Logger) Debug(message string) {
	if l.debug {
		fmt.Println(color.MagentaString(message))
	}
}

// Error logs a message to the console as an error.
func (l *Logger) Error(message string) {
	fmt.Println(color.RedString(message))
}

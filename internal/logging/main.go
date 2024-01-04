package logging

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// Logger is a struct that represents a logger.
type Logger struct {
	DebugFlag bool
}

// NewLogger creates a new Logger instance.
func NewLogger(debug bool) *Logger {
	return &Logger{
		DebugFlag: debug,
	}
}

// Log logs a message to the console.
func (l *Logger) Log(message string) {
	fmt.Println(color.WhiteString(message))
}

// Debug logs a message to the console if debug is enabled.
func (l *Logger) Debug(message string) {
	if l.DebugFlag {
		fmt.Println(color.MagentaString(message))
	}
}

// Error logs a message to the console as an error.
func (l *Logger) Error(message string) {
	fmt.Println(color.RedString(message))
}

// Fatal logs a message to the console as an error and exits the program.
func (l *Logger) Fatal(message string) {
	fmt.Println(color.RedString(message))
	os.Exit(1)
}

// FatalErr logs an error to the console as an error and exits the program.
func (l *Logger) FatalErr(err error) {
	fmt.Println(color.RedString(err.Error()))
	os.Exit(1)
}

// Success logs a message to the console as a success.
func (l *Logger) Success(message string) {
	fmt.Println(color.GreenString(message))
}

// Warning logs a message to the console as a warning.
func (l *Logger) Warning(message string) {
	fmt.Println(color.YellowString(message))
}

// Info logs a message to the console as an info.
func (l *Logger) Info(message string) {
	fmt.Println(color.BlueString(message))
}

// Logf logs a formatted message to the console.
func (l *Logger) Logf(format string, a ...interface{}) {
	fmt.Printf(color.WhiteString(format), a...)
}

// Debugf logs a formatted message to the console if debug is enabled.
func (l *Logger) Debugf(format string, a ...interface{}) {
	if l.DebugFlag {
		fmt.Printf(color.MagentaString(format), a...)
	}
}

// Errorf logs a formatted message to the console as an error.
func (l *Logger) Errorf(format string, a ...interface{}) {
	fmt.Printf(color.RedString(format), a...)
}

// Fatalf logs a formatted message to the console as an error and exits the program.
func (l *Logger) Fatalf(format string, a ...interface{}) {
	fmt.Printf(color.RedString(format), a...)
	os.Exit(1)
}

// Successf logs a formatted message to the console as a success.
func (l *Logger) Successf(format string, a ...interface{}) {
	fmt.Printf(color.GreenString(format), a...)
}

// Warningf logs a formatted message to the console as a warning.
func (l *Logger) Warningf(format string, a ...interface{}) {
	fmt.Printf(color.YellowString(format), a...)
}

// Infof logs a formatted message to the console as an info.
func (l *Logger) Infof(format string, a ...interface{}) {
	fmt.Printf(color.BlueString(format), a...)
}

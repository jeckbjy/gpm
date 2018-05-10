package gpm

import (
	"fmt"
	"os"
	"sync"
)

// These contanstants map to color codes for shell scripts making them
// human readable.
const (
	Blue   = "0;34"
	Red    = "0;31"
	Green  = "0;32"
	Yellow = "0;33"
	Cyan   = "0;36"
	Pink   = "1;35"
)

const (
	zLogInfo = iota
	zLogDebug
	zLogWarn
	zLogError
	zLogFatal
)

var (
	zLogColor = []string{Green, "", Yellow, Red, Red}
	zLogName  = []string{"INFO", "DEBUG", "WARN", "ERROR", "FATAL"}
)

func NewLogger() *Logger {
	l := &Logger{Quiet: false, Debuging: false, NoColor: false, PanicOnDie: false}
	return l
}

// Logger system
type Logger struct {
	sync.Mutex
	Quiet      bool
	Debuging   bool
	NoColor    bool
	PanicOnDie bool
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.Output(zLogInfo, msg, args...)
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.Output(zLogDebug, msg, args...)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.Output(zLogWarn, msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.Output(zLogError, msg, args...)
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.Output(zLogFatal, msg, args...)
}

func (l *Logger) Output(level int, msg string, args ...interface{}) {
	l.Lock()
	defer l.Unlock()
	name := zLogName[level]
	fmt.Printf("[%s]\t", name)
	fmt.Printf(msg, args...)
	fmt.Println("")
}

// Print prints exactly the string given.
//
// It prints to Stdout.
func (l *Logger) Print(msg string) {
	l.Lock()
	defer l.Unlock()
	fmt.Fprint(os.Stdout, msg)
}

// Puts formats a message and then prints to Stdout.
//
// It does not prefix the message, does not color it, or otherwise decorate it.
//
// It does add a line feed.
func (l *Logger) Puts(msg string, args ...interface{}) {
	l.Lock()
	defer l.Unlock()

	fmt.Fprintf(os.Stdout, msg, args...)
	fmt.Fprintln(os.Stdout)
}

func (l *Logger) Die(msg string, args ...interface{}) {
	l.Exit(0, msg, args...)
}

func (l *Logger) Exit(code int, msg string, args ...interface{}) {
	l.Error(msg, args...)

	if l.PanicOnDie {
		panic("trapped")
	}

	os.Exit(code)
}

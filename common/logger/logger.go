package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

type lvl int32

const (
	TRACE lvl = iota
	DEBUG
	INFO
	ERROR
	PANIC
	DISABLED
)

var (
	globalLevel = new(int32)
)

func (l lvl) String() string {
	switch l {
	case TRACE:
		return "Trace"
	case DEBUG:
		return "Debug"
	case INFO:
		return "Info"
	case ERROR:
		return "Error"
	case PANIC:
		return "Panic"
	case DISABLED:
		return ""
	default:
		return ""
	}
}

func setLevel(l lvl) {
	atomic.StoreInt32(globalLevel, int32(l))
}

func getLevel() lvl {
	return lvl(atomic.LoadInt32(globalLevel))
}

type logger struct {
	w io.Writer
}

func (l *logger) setLoggerOutput(w io.Writer) {
	l.w = w
}

func (l *logger) print(lv lvl, msg string) {
	if lv < getLevel() {
		return
	}
	tm := time.Now()
	s := fmt.Sprintf("[%s] %s: %s\n", tm.Format(time.RFC3339Nano), lv.String(), msg)
	l.w.Write([]byte(s))
}

var log = newLogger(os.Stderr)

func SetLogger(w io.Writer) {
	log.setLoggerOutput(w)
}

func SetLevel(l string) {
	switch strings.ToUpper(l) {
	case "TRACE":
		setLevel(TRACE)
	case "DEBUG":
		setLevel(DEBUG)
	case "INFO":
		setLevel(INFO)
	case "ERROR":
		setLevel(ERROR)
	case "PANIC":
		setLevel(PANIC)
	case "DISABLED":
		setLevel(DISABLED)
	default:
		setLevel(INFO)
	}
}

func Tracef(f string, v ...interface{}) {
	log.print(TRACE, fmt.Sprintf(f, v...))
}

func Debugf(f string, v ...interface{}) {
	log.print(DEBUG, fmt.Sprintf(f, v...))
}

func Infof(f string, v ...interface{}) {
	log.print(INFO, fmt.Sprintf(f, v...))
}

func Errorf(f string, v ...interface{}) {
	log.print(ERROR, fmt.Sprintf(f, v...))
}

func Panicf(f string, v ...interface{}) {
	log.print(PANIC, fmt.Sprintf(f, v...))
	panic("PANIC!")
}

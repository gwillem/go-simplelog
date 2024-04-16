package log

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/fatih/color"
)

type (
	Level             int
	prefixerInterface interface {
		Prefix() string
	}
	defaultPrefixer struct{}
)

const (
	LevelDebug Level = iota
	LevelTask
	LevelWarn
	LevelAlert
	LevelError
)

var (
	Yellow     = color.New(color.FgYellow, color.Bold).SprintFunc()
	Red        = color.New(color.FgHiRed, color.Bold).SprintFunc()
	Purple     = color.New(color.FgMagenta, color.Bold).SprintFunc()
	Green      = color.New(color.FgGreen, color.Bold).SprintFunc()
	WhiteOnRed = color.New(color.FgWhite, color.BgRed, color.Bold).SprintFunc()
	Dark       = color.New(color.FgHiBlack).SprintFunc()
	BoldWhite  = color.New(color.FgHiWhite, color.Bold).SprintFunc()

	channel  io.Writer = os.Stderr
	logLevel           = LevelDebug

	writeLock sync.Mutex
	prefixer  prefixerInterface = defaultPrefixer{}
)

func (p defaultPrefixer) Prefix() string {
	return ""
}

func SetPrefixer(p prefixerInterface) {
	prefixer = p
}

func Silence(new bool) bool {
	prev := channel == io.Discard

	if new {
		channel = io.Discard
	} else {
		channel = os.Stderr
	}

	return prev
}

func SetLevel(l Level) {
	logLevel = l
}

func init() {
	// color.NoColor = false // Override terminal detection
}

func Debug(arg ...interface{}) {
	if logLevel <= LevelDebug {
		_print(Dark("   "), arg...)
	}
}

func Task(arg ...interface{}) {
	if logLevel <= LevelTask {
		_print(Yellow(">>>"), arg...)
	}
}

func Warn(arg ...interface{}) {
	if logLevel <= LevelWarn {
		_print(Red("!!!"), arg...)
	}
}

func Alert(arg ...interface{}) {
	if logLevel <= LevelAlert {
		_print(Purple(" ! "), arg...)
	}
}

func Ok(arg ...interface{}) {
	if logLevel <= LevelTask {
		_print(Green(" ✔ "), arg...)
	}
}

func Progress(arg ...interface{}) {
	if logLevel <= LevelTask {
		_print(" - ", arg...)
	}
}

func Fatal(arg ...interface{}) {
	_print(WhiteOnRed("XXX"), arg...)
	os.Exit(1)
}

func Check(e error, msg ...interface{}) {
	if e != nil {
		if len(msg) > 0 {
			_print("ERR", msg...)
		}
		Fatal(e.Error())
	}
}

func Error(e error) {
	Fatal("Fatal error:", e.Error())
}

func _print(prefix string, arg ...interface{}) {
	writeLock.Lock()
	defer writeLock.Unlock()

	fmt.Fprintf(channel, "%s %s %s", prefixer.Prefix(), prefix, fmt.Sprintln(arg...))
}

package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/fatih/color"
)

type (
	Level             int
	prefixerInterface interface {
		Prefix() string
		PrefixDebug() string
		PrefixTask() string
		PrefixWarn() string
		PrefixAlert() string
		PrefixOk() string
		PrefixProgress() string
		PrefixFatal() string
		PrefixError() string
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

func (p defaultPrefixer) PrefixDebug() string {
	return Dark("   ")
}

func (p defaultPrefixer) PrefixTask() string {
	return Yellow(">>>")
}

func (p defaultPrefixer) PrefixWarn() string {
	return Red("!!!")
}

func (p defaultPrefixer) PrefixAlert() string {
	return Purple(" ! ")
}

func (p defaultPrefixer) PrefixOk() string {
	return Green(" ✔ ")
}

func (p defaultPrefixer) PrefixProgress() string {
	return " - "
}

func (p defaultPrefixer) PrefixFatal() string {
	return WhiteOnRed("XXX")
}

func (p defaultPrefixer) PrefixError() string {
	return Red("ERR")
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

func IsSilenced() bool {
	return channel == io.Discard
}

func SetLevel(l Level) {
	logLevel = l
}

func GetLevel() Level {
	return logLevel
}

func Debug(arg ...any) {
	if logLevel <= LevelDebug {
		_print(prefixer.PrefixDebug(), arg...)
	}
}

func Debugf(format string, arg ...any) {
	if logLevel <= LevelDebug {
		_print(prefixer.PrefixDebug(), fmt.Sprintf(format, arg...))
	}
}

func Task(arg ...any) {
	if logLevel <= LevelTask {
		_print(prefixer.PrefixTask(), arg...)
	}
}

func Taskf(format string, arg ...any) {
	if logLevel <= LevelTask {
		_print(prefixer.PrefixTask(), fmt.Sprintf(format, arg...))
	}
}

func Warn(arg ...any) {
	if logLevel <= LevelWarn {
		_print(prefixer.PrefixWarn(), arg...)
	}
}

func Warnf(format string, arg ...any) {
	if logLevel <= LevelWarn {
		_print(prefixer.PrefixWarn(), fmt.Sprintf(format, arg...))
	}
}

func Alert(arg ...any) {
	if logLevel <= LevelAlert {
		_print(prefixer.PrefixAlert(), arg...)
	}
}

func Alertf(format string, arg ...any) {
	if logLevel <= LevelAlert {
		_print(prefixer.PrefixAlert(), fmt.Sprintf(format, arg...))
	}
}

func Ok(arg ...any) {
	if logLevel <= LevelTask {
		_print(prefixer.PrefixOk(), arg...)
	}
}

func Okf(format string, arg ...any) {
	if logLevel <= LevelTask {
		_print(prefixer.PrefixOk(), fmt.Sprintf(format, arg...))
	}
}

func Progress(arg ...any) {
	if logLevel <= LevelTask {
		_print(prefixer.PrefixProgress(), arg...)
	}
}

func Progressf(format string, arg ...any) {
	if logLevel <= LevelTask {
		_print(prefixer.PrefixProgress(), fmt.Sprintf(format, arg...))
	}
}

func Fatal(arg ...any) {
	_print(prefixer.PrefixFatal(), arg...)
	os.Exit(1)
}

func Fatalf(format string, arg ...any) {
	Fatal(fmt.Sprintf(format, arg...))
}

func Check(e error, msg ...any) {
	if e != nil {
		if len(msg) > 0 {
			_print(prefixer.PrefixError(), msg...)
		}
		Fatal(e.Error())
	}
}

func Checkf(e error, format string, arg ...any) {
	if e != nil {
		if len(arg) > 0 {
			_print(prefixer.PrefixError(), fmt.Sprintf(format, arg...))
		}
		Fatal(e.Error())
	}
}

func Must[T any](v T, err error) T {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		Fatal(fmt.Sprintf("%s:%d: %v", lastTwoSegments(file), line, err))
	}
	return v
}

func Error(e error) {
	Fatal("Fatal error:", e.Error())
}

func Errorf(format string, arg ...any) {
	Fatal(fmt.Errorf(format, arg...))
}

func _print(prefix string, arg ...any) {
	writeLock.Lock()
	defer writeLock.Unlock()

	fmt.Fprintf(channel, "%s %s %s", prefixer.Prefix(), prefix, fmt.Sprintln(arg...))
}

func lastTwoSegments(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return path
	}
	return strings.Join(parts[len(parts)-2:], "/")
}

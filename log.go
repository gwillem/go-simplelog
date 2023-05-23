package log

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

var (
	Yellow     = color.New(color.FgYellow, color.Bold).SprintFunc()
	Red        = color.New(color.FgHiRed, color.Bold).SprintFunc()
	Purple     = color.New(color.FgMagenta, color.Bold).SprintFunc()
	Green      = color.New(color.FgGreen, color.Bold).SprintFunc()
	WhiteOnRed = color.New(color.FgWhite, color.BgRed, color.Bold).SprintFunc()
	Dark       = color.New(color.FgHiBlack).SprintFunc()
	BoldWhite  = color.New(color.FgHiWhite, color.Bold).SprintFunc()

	channel io.Writer = os.Stderr
)

func Silence(new bool) bool {

	prev := channel == io.Discard

	if new {
		channel = io.Discard
	} else {
		channel = os.Stderr
	}

	return prev
}

func init() {
	// color.NoColor = false // Override terminal detection
}

func Debug(arg ...interface{}) {
	print(Dark("   "), arg...)
}

func Task(arg ...interface{}) {
	print(Yellow(">>>"), arg...)
}

func Warn(arg ...interface{}) {
	print(Red("!!!"), arg...)
}

func Alert(arg ...interface{}) {
	print(Purple(" ! "), arg...)
}

func Ok(arg ...interface{}) {
	print(Green(" ✔ "), arg...)
}

func Progress(arg ...interface{}) {
	print(" - ", arg...)
}

func Fatal(arg ...interface{}) {
	print(WhiteOnRed("XXX"), arg...)
	os.Exit(1)
}

func Check(e error) {
	if e != nil {
		Fatal(e.Error())
	}
}

func Error(e error) {
	Fatal("Fatal error:", e.Error())
}

func print(prefix string, arg ...interface{}) {
	fmt.Fprint(channel, prefix+" ")
	fmt.Fprintln(channel, arg...)
}

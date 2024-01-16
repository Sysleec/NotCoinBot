package clicker

import (
	"log"
	"os"

	"github.com/fatih/color"
)

var (
	InfoLogger    = NewLogger("INFO: ", color.New(color.FgWhite).SprintFunc())
	SuccessLogger = NewLogger("SUCCESS: ", color.New(color.FgGreen).SprintFunc())
	WarningLogger = NewLogger("WARNING: ", color.New(color.FgYellow).SprintFunc())
	ErrorLogger   = NewLogger("ERROR: ", color.New(color.FgRed).SprintFunc())
	DebugLogger   = NewLogger("DEBUG: ", color.New(color.FgHiBlack).SprintFunc())
)

func (cw *ColorableWriter) Write(p []byte) (n int, err error) {
	cw.Console.Write([]byte(cw.Color(string(p))))
	return cw.File.Write(p)
}

func NewLogger(prefix string, colorFunc func(a ...interface{}) string) *log.Logger {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	cw := &ColorableWriter{
		Console: os.Stdout,
		File:    file,
		Prefix:  prefix,
		Color:   colorFunc,
	}

	return log.New(cw, prefix, log.Ldate|log.Ltime|log.Lshortfile)
}

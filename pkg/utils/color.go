package utils

import (
	"fmt"

	"github.com/mgutz/ansi"
)

// Outputs ANSI color if stdout is a tty
var (
	Magenta          = makeColorFunc("magenta")
	MagentaFormatted = makeColorFuncFormatted("magenta")
	Cyan             = makeColorFunc("cyan")
	CyanFormatted    = makeColorFuncFormatted("cyan")
	Red              = makeColorFunc("red")
	RedFormatted     = makeColorFuncFormatted("red")
	Yellow           = makeColorFunc("yellow")
	YellowFormatted  = makeColorFuncFormatted("yellow")
	Blue             = makeColorFunc("blue")
	BlueFormatted    = makeColorFuncFormatted("blue")
	Green            = makeColorFunc("green")
	GreenFormatted   = makeColorFuncFormatted("green")
	Gray             = makeColorFunc("black+h")
	GrayFormatted    = makeColorFuncFormatted("black+h")
	Bold             = makeColorFunc("default+b")
	BoldFormatted    = makeColorFuncFormatted("default+b")
)

func makeColorFunc(color string) func(string) string {
	cf := ansi.ColorFunc(color)
	return func(arg string) string {
		return cf(arg)
	}
}

func makeColorFuncFormatted(color string) func(string, string) string {
	cf := ansi.ColorFunc(color)
	return func(arg string, format string) string {
		if format != "" {
			return cf(fmt.Sprintf(format, arg))
		}
		return cf(arg)
	}
}

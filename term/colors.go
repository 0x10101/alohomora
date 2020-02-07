package term

import (
	"fmt"
	"os"
)

var (
	AnsiBold = "\033[1m"
	AnsiDim  = "\033[2m"

	AnsiRed     = "\033[31m"
	AnsiGreen   = "\033[32m"
	AnsiYellow  = "\033[33m"
	AnsiBlue    = "\033[34m"
	AnsiMagenta = "\033[35m"
	AnsiCyan    = "\033[36m"
	AnsiBlack   = "\033[30m"
	AnsiWhite   = "\033[97m"

	AnsiBgDarkGray = "\033[100m"
	AnsiBgRed      = "\033[41m"
	AnsiBgGreen    = "\033[42m"
	AnsiBgYellow   = "\033[43m"
	AnsiBgMagenta  = "\033[45m"
	AnsiBgBlue     = "\033[104m"
	AnsiBgCyan     = "\033[46m"

	AnsiReset = "\033[0m"
)

func Colors() {
	AnsiBold = "\033[1m"
	AnsiDim = "\033[2m"

	AnsiRed = "\033[31m"
	AnsiGreen = "\033[32m"
	AnsiYellow = "\033[33m"
	AnsiBlue = "\033[34m"
	AnsiMagenta = "\033[35m"
	AnsiBlack = "\033[30m"
	AnsiWhite = "\033[97m"
	AnsiCyan = "\033[36m"

	AnsiBgDarkGray = "\033[100m"
	AnsiBgRed = "\033[41m"
	AnsiBgGreen = "\033[42m"
	AnsiBgYellow = "\033[43m"
	AnsiBgMagenta = "\033[45m"
	AnsiBgCyan = "\033[46m"
	AnsiBgBlue = "\033[104m"

	AnsiReset = "\033[0m"
}

func NoColors() {
	AnsiBold = ""
	AnsiDim = ""

	AnsiRed = ""
	AnsiGreen = ""
	AnsiYellow = ""
	AnsiBlue = ""
	AnsiMagenta = ""
	AnsiBlack = ""
	AnsiWhite = ""
	AnsiCyan = ""

	AnsiBgDarkGray = ""
	AnsiBgRed = ""
	AnsiBgGreen = ""
	AnsiBgYellow = ""
	AnsiBgMagenta = ""
	AnsiBgBlue = ""
	AnsiBgCyan = ""

	AnsiReset = ""
}

func wrap(ansi, value string) string {
	return fmt.Sprintf("%s%s%s", ansi, value, AnsiReset)
}

func BgCyan(value string) string {
	return wrap(AnsiBgCyan, value)
}

func BgGray(value string) string {
	return wrap(AnsiBgDarkGray, value)
}

func BgRed(value string) string {
	return wrap(AnsiBgRed, value)
}

func BgGreen(value string) string {
	return wrap(AnsiBgGreen, value)
}

func BgYellow(value string) string {
	return wrap(AnsiBgYellow, value)
}

func BgBlue(value string) string {
	return wrap(AnsiBgBlue, value)
}

func BgMagenta(value string) string {
	return wrap(AnsiBgMagenta, value)
}

func Bold(value string) string {
	return wrap(AnsiBold, value)
}

func Cyan(value string) string {
	return wrap(AnsiCyan, value)
}

func Red(value string) string {
	return wrap(AnsiRed, value)
}

func BrightRed(value string) string {
	return Red(Bold(value))
}

func Green(value string) string {
	return wrap(AnsiGreen, value)
}

func BrightGreen(value string) string {
	return Green(Bold(value))
}

func Yellow(value string) string {
	return wrap(AnsiYellow, value)
}

func Blue(value string) string {
	return wrap(AnsiBlue, value)
}

func BrightBlue(value string) string {
	return Blue(Bold(value))
}

func Black(value string) string {
	return wrap(AnsiBlack, value)
}

func White(value string) string {
	return wrap(AnsiWhite, value)
}

func Magenta(value string) string {
	return wrap(AnsiMagenta, value)
}

func BrightMagenta(value string) string {
	return Magenta(Bold(value))
}

func BrightCyan(value string) string {
	return Cyan(Bold(value))
}

func BrightYellow(value string) string {
	return Yellow(Bold(value))
}

func Dim(value string) string {
	return wrap(AnsiDim, value)
}

func DimBlue(value string) string {
	return Dim(Blue(value))
}

func DimGreen(value string) string {
	return Dim(Green(value))
}

func DimYellow(value string) string {
	return Dim(Yellow(value))
}

func DimMagenta(value string) string {
	return Dim(Magenta(value))
}

func DimRed(value string) string {
	return Dim(Red(value))
}

func DimCyan(value string) string {
	return Dim(Cyan(value))
}

func Supported() bool {
	if term := os.Getenv("TERM"); term == "" {
		return false
	} else if term == "dumb" {
		return false
	}
	return true
}

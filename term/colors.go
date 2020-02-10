package term

import (
	"fmt"
	"os"
)

var (
	// AnsiBold is the escape sequence for bold/emphasized text
	AnsiBold = "\033[1m"

	// AnsiDim is the escape sequence for dimmed text
	AnsiDim = "\033[2m"

	// AnsiRed is the escape sequence for red text
	AnsiRed = "\033[31m"

	// AnsiGreen is the escape sequence for green text
	AnsiGreen = "\033[32m"

	// AnsiYellow is the escape sequence for yellow text
	AnsiYellow = "\033[33m"

	// AnsiBlue is the escape sequence for blue text
	AnsiBlue = "\033[34m"

	// AnsiMagenta is the escape sequence for magenta text
	AnsiMagenta = "\033[35m"

	// AnsiCyan is the escape sequence for cyan text
	AnsiCyan = "\033[36m"

	// AnsiBlack is the escape sequence for black text
	AnsiBlack = "\033[30m"

	// AnsiWhite is the escape sequence for white text
	AnsiWhite = "\033[97m"

	// AnsiBgDarkGray is the escape sequence for a dark gray background
	AnsiBgDarkGray = "\033[100m"

	// AnsiBgRed is the escape sequence for red background
	AnsiBgRed = "\033[41m"

	// AnsiBgGreen is the escape sequence for green background
	AnsiBgGreen = "\033[42m"

	// AnsiBgYellow is the escape sequence for yellow background
	AnsiBgYellow = "\033[43m"

	// AnsiBgMagenta is the escape sequence for magenta background
	AnsiBgMagenta = "\033[45m"

	// AnsiBgBlue is the escape sequence for blue background
	AnsiBgBlue = "\033[104m"

	// AnsiBgCyan is the escape sequence for cyan background
	AnsiBgCyan = "\033[46m"

	// AnsiReset resets formatting to default
	AnsiReset = "\033[0m"
)

// Colors enables color escape sequences on supported terminals
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

// NoColors disables colors on terminals that don't understand ANSI escape sequences
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

// BgCyan formats a given string by prepending the AnsiBgCyan escape sequence and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func BgCyan(value string) string {
	return wrap(AnsiBgCyan, value)
}

// BgGray formats a given string by prepending the AnsiBgDarkGray escape sequence and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func BgGray(value string) string {
	return wrap(AnsiBgDarkGray, value)
}

// BgRed formats a given string by prepending the AnsiBgRed escape sequence and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func BgRed(value string) string {
	return wrap(AnsiBgRed, value)
}

// BgGreen formats a given string by prepending the AnsiBgGreen escape sequence and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func BgGreen(value string) string {
	return wrap(AnsiBgGreen, value)
}

// BgYellow formats a given string by prepending the AnsiBgYellow escape sequence and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func BgYellow(value string) string {
	return wrap(AnsiBgYellow, value)
}

// BgBlue formats a given string by prepending the AnsiBgBlue escape sequence and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func BgBlue(value string) string {
	return wrap(AnsiBgBlue, value)
}

// BgMagenta formats a given string by prepending the AnsiBgMagenta escape sequence and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func BgMagenta(value string) string {
	return wrap(AnsiBgMagenta, value)
}

// Bold formats a given string by prepending the AnsiBold escape sequence and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func Bold(value string) string {
	return wrap(AnsiBold, value)
}

// Cyan formats a given string by prepending the AnsiCyan escape sequence and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func Cyan(value string) string {
	return wrap(AnsiCyan, value)
}

// Red formats a given string by prepending the AnsiRed escape sequence and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func Red(value string) string {
	return wrap(AnsiRed, value)
}

// BrightRed formats a given string by prepending the AnsiBold and AnsiRed escape sequences and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func BrightRed(value string) string {
	return Red(Bold(value))
}

// Green formats a given string by prepending the AnsiGreen escape sequence and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func Green(value string) string {
	return wrap(AnsiGreen, value)
}

// BrightGreen formats a given string by prepending the AnsiGreen and AnsiBold escape sequences and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func BrightGreen(value string) string {
	return Green(Bold(value))
}

// Yellow formats a given string by prepending the AnsiYellow escape sequence and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func Yellow(value string) string {
	return wrap(AnsiYellow, value)
}

// Blue formats a given string by prepending the AnsiBlue escape sequence and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func Blue(value string) string {
	return wrap(AnsiBlue, value)
}

// BrightBlue formats a given string by prepending the AnsiBlue and AnsiBold escape sequences and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func BrightBlue(value string) string {
	return Blue(Bold(value))
}

// Black formats a given string by prepending the AnsiBlack escape sequence and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func Black(value string) string {
	return wrap(AnsiBlack, value)
}

// White formats a given string by prepending the AnsiWhite escape sequence and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func White(value string) string {
	return wrap(AnsiWhite, value)
}

// Magenta formats a given string by prepending the AnsiMagenta escape sequence and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func Magenta(value string) string {
	return wrap(AnsiMagenta, value)
}

// BrightMagenta formats a given string by prepending the AnsiMagenta and AnsiBold escape sequences and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func BrightMagenta(value string) string {
	return Magenta(Bold(value))
}

// BrightCyan formats a given string by prepending the AnsiCyan and AnsiBold escape sequences and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func BrightCyan(value string) string {
	return Cyan(Bold(value))
}

// BrightYellow formats a given string by prepending the AnsiYellow and AnsiBold escape sequences and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func BrightYellow(value string) string {
	return Yellow(Bold(value))
}

// Dim formats a given string by prepending the AnsiDim escape sequence and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func Dim(value string) string {
	return wrap(AnsiDim, value)
}

// DimBlue formats a given string by prepending the AnsiBim and AnsiBlue escape sequences and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func DimBlue(value string) string {
	return Dim(Blue(value))
}

// DimGreen formats a given string by prepending the AnsiDim and AnsiGreen escape sequences and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func DimGreen(value string) string {
	return Dim(Green(value))
}

// DimYellow formats a given string by prepending the AnsiDim and AnsiYellow escape sequence and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func DimYellow(value string) string {
	return Dim(Yellow(value))
}

// DimMagenta formats a given string by prepending the AnsiDim and AnsiMagenta escape sequences and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func DimMagenta(value string) string {
	return Dim(Magenta(value))
}

// DimRed formats a given string by prepending the AnsiDim and AnsiRed escape sequences and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func DimRed(value string) string {
	return Dim(Red(value))
}

// DimCyan formats a given string by prepending the AnsiDim and AnsiCyan escape sequences and postfixing the AnsiReset escape sequence.
// The formatted string is returned and can be printed to supported terminals. If colors are disabled, the escape sequences
// are replaced by empty strings, so that this function can safely be used on non-supporting terminals.
func DimCyan(value string) string {
	return Dim(Cyan(value))
}

// Supported determines whether or not the current terminal supports ANSI escape sequences.
// If they are supported, true is returned, otherwise false.
func Supported() bool {
	if term := os.Getenv("TERM"); term == "" {
		return false
	} else if term == "dumb" {
		return false
	}
	return true
}

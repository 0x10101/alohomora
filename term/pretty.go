package term

import (
	"fmt"
)

// prettyMessage helps formatting text for the convenience functions defined in this
// file. It prefixes a message with square brackets containing a (colored, if supported)
// marker like [*] (info), [!] (warning) and so on.
func prettyMessage(format, prefix string, data ...interface{}) {
	fmt.Printf("%s %s", prefix, fmt.Sprintf(format, data...))
}

// Info conveniently wraps a message in the way alohomora is printing text to stdout.
// The message will be formatted this way: [*] <MSG>. The * will be cyan (if supported).
func Info(format string, data ...interface{}) {
	prettyMessage(format, fmt.Sprintf("[%s]", Cyan("*")), data...)
}

// Warn conveniently wraps a message in the way alohomora is printing text to stdout.
// The message will be formatted this way: [!] <MSG>. The ! will be yellow (if supported).
func Warn(format string, data ...interface{}) {
	prettyMessage(format, fmt.Sprintf("[%s]", Yellow("!")), data...)
}

// Problem conveniently wraps a message in the way alohomora is printing text to stdout.
// The message will be formatted this way: [-] <MSG>. The - will be red (if supported).
func Problem(format string, data ...interface{}) {
	prettyMessage(format, fmt.Sprintf("[%s]", Red("-")), data...)
}

// Error conveniently wraps a message in the way alohomora is printing text to stdout.
// The message will be formatted this way: [x] <MSG>. The x will be red (if supported).
func Error(format string, data ...interface{}) {
	//prettyMessage(format, fmt.Sprintf("[%s]", Red("x")), data...)
	fmt.Printf("%s %s%s%s%s", fmt.Sprintf("[%s]", Red("x")), AnsiBold, AnsiRed, fmt.Sprintf(format, data...), AnsiReset)
}

// Success conveniently wraps a message in the way alohomora is printing text to stdout.
// The message will be formatted this way: [+] <MSG>. The + will be green (if supported).
func Success(format string, data ...interface{}) {
	prettyMessage(format, fmt.Sprintf("[%s]", Green("+")), data...)
}

// LabelMagenta prints a magenta label with white foreground to stdout.
func LabelMagenta(data string) string {
	return BgMagenta(Bold(White(fmt.Sprintf(" %s ", data))))
}

// LabelGreen prints a green label with white foreground to stdout.
func LabelGreen(data string) string {
	return BgGreen(Bold(White(fmt.Sprintf(" %s ", data))))
}

// LabelBlue prints a blue label with white foreground to stdout.
func LabelBlue(data string) string {
	return BgBlue(Bold(White(fmt.Sprintf(" %s ", data))))
}

// LabelRed prints a red label with white foreground to stdout.
func LabelRed(data string) string {
	return BgRed(Bold(White(fmt.Sprintf(" %s ", data))))
}

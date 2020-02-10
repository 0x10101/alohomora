package term

import (
	"fmt"
)

func prettyMessage(format, prefix string, data ...interface{}) {
	fmt.Printf("%s %s", prefix, fmt.Sprintf(format, data...))
}

func Info(format string, data ...interface{}) {
	prettyMessage(format, fmt.Sprintf("[%s]", Cyan("*")), data...)
}

func Warn(format string, data ...interface{}) {
	prettyMessage(format, fmt.Sprintf("[%s]", Yellow("!")), data...)
}

func Problem(format string, data ...interface{}) {
	prettyMessage(format, fmt.Sprintf("[%s]", Red("-")), data...)
}

func Error(format string, data ...interface{}) {
	//prettyMessage(format, fmt.Sprintf("[%s]", Red("x")), data...)
	fmt.Printf("%s %s%s%s%s", fmt.Sprintf("[%s]", Red("x")), AnsiBold, AnsiRed, fmt.Sprintf(format, data...), AnsiReset)
}

func Success(format string, data ...interface{}) {
	prettyMessage(format, fmt.Sprintf("[%s]", Green("+")), data...)
}

func LabelMagenta(data string) string {
	return BgMagenta(Bold(White(fmt.Sprintf(" %s ", data))))
}

func LabelGreen(data string) string {
	return BgGreen(Bold(White(fmt.Sprintf(" %s ", data))))
}

func LabelBlue(data string) string {
	return BgBlue(Bold(White(fmt.Sprintf(" %s ", data))))
}

func LabelRed(data string) string {
	return BgRed(Bold(White(fmt.Sprintf(" %s ", data))))
}

package term

import (
	"fmt"
	"strings"
)

type Table struct {
	cols           int
	raw            [][]string
	widths         []int
	Separator      string
	MaxColWidth    int
	maxActualWidth int
}

func NewTable() *Table {
	t := new(Table)
	t.raw = make([][]string, 0)
	t.widths = make([]int, 0)
	t.Separator = "\t"
	return t
}

func (t *Table) AddRow(data ...string) {
	t.raw = append(t.raw, data)
	if len(data) > t.cols {
		t.cols = len(data)
	}
}

func wFmt(width int) string {
	return fmt.Sprintf("%%-%ds", width)
}

func (t *Table) spaces(colIndex int) string {
	var i int
	var sb strings.Builder

	for i < colIndex {
		j := t.widths[i]
		for j > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(t.Separator)
	}

	return sb.String()
}

func (t *Table) Format() string {

	var sb strings.Builder

	var i int
	t.widths = make([]int, t.cols)

	for _, arr := range t.raw {
		i = 0
		for _, str := range arr {
			l := len(str)
			if t.maxActualWidth < l {
				t.maxActualWidth = l
			}
			if l > t.widths[i] {
				t.widths[i] = l
			}
			i++
		}
	}

	fmt.Println(t.widths)
	for _, arr := range t.raw {
		i = 0

		for colIndex, str := range arr {
			w := t.widths[i]
			i++

			if len(str) > t.MaxColWidth && t.MaxColWidth > 0 {
				// Split the string in parts
				start := 0
				end := t.MaxColWidth

				for {
					fmt.Printf("Slicing %d char string from %d to %d\n", len(str), start, end)
					part := str[start:end]
					start = end + 1
					if start > len(str) {
						start = len(str)
					}
					end = end + t.MaxColWidth + 1
					if end > len(str) {
						end = len(str)
					}
					fStr := fmt.Sprintf("%s%s", wFmt(w), t.Separator)
					sb.WriteString(t.spaces(colIndex))
					sb.WriteString(fmt.Sprintf(fStr, part))
					sb.WriteString("\n")
					if start == end {
						break
					}
				}
			} else {
				fStr := fmt.Sprintf("%s%s", wFmt(w), t.Separator)
				sb.WriteString(fmt.Sprintf(fStr, str))
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()

}

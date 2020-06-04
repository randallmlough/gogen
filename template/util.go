package template

import (
	"bytes"
	"strings"
)

func beginRegion(buf bytes.Buffer, region string) {
	buf.WriteString("\n// region    " + center(70, "*", " "+region+" ") + "\n")
}

func endRegion(buf bytes.Buffer, region string) {
	buf.WriteString("\n// endregion " + center(70, "*", " "+region+" ") + "\n")
}

func center(width int, pad string, s string) string {
	if len(s)+2 > width {
		return s
	}
	lpad := (width - len(s)) / 2
	rpad := width - (lpad + len(s))
	return strings.Repeat(pad, lpad) + s + strings.Repeat(pad, rpad)
}

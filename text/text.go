package text

import (
	"fmt"

	"github.com/mattn/go-runewidth"
)

var rw *runewidth.Condition

func init() {
	rw = runewidth.NewCondition()
	rw.EastAsianWidth = true
	rw.StrictEmojiNeutral = true
}

func PadRight(v any, len int) string {
	return rw.FillRight(fmt.Sprint(v), len)
}

func PadLeft(v any, len int) string {
	return rw.FillLeft(fmt.Sprint(v), len)
}

func StringWidth(v any) int {
	return rw.StringWidth(fmt.Sprint(v))
}

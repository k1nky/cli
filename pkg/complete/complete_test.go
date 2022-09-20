package complete

import (
	"testing"

	"github.com/chzyer/readline"
)

func isRuneSliceEqual(a [][]rune, b [][]rune) bool {
	if len(a) != len(b) {
		return false
	}
	runes := readline.Runes{}
	for _, va := range a {
		found := false
		for _, vb := range b {
			if runes.Equal(va, vb) {
				found = true
				continue
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func TestDo(t *testing.T) {
	c := CompleterItem("", nil, CompleterItem("cmd1", nil, CompleterItem("cmd2", []string{"aaa", "bbb"}), CompleterItem("cmd3", nil)))
	cases := []struct {
		line            []rune
		expectedOffset  int
		expectedNewLine [][]rune
	}{
		{[]rune("cmd1 cmd2 aaa=1234 bbb=12345 "), 0, [][]rune{[]rune("aaa="), []rune("bbb=")}},
		{[]rune("cmd1 cmd2 "), 0, [][]rune{[]rune("aaa="), []rune("bbb=")}},
		{[]rune("cmd1 cmd2 a"), 1, [][]rune{[]rune("aa=")}},
		{[]rune("cmd1 cmd2 aaa="), 0, [][]rune{}},
		{[]rune("cmd1 cmd2 aaa=value"), 9, [][]rune{[]rune(" ")}},
		{[]rune("cmd1 cmd2  aaa=value bbb=value "), 0, [][]rune{[]rune("aaa="), []rune("bbb=")}},
		{[]rune("cmd1 cmd"), 3, [][]rune{[]rune("2 "), []rune("3 ")}},
	}
	for k, v := range cases {
		newLine, offset := c.Do(v.line, len(v.line))
		if offset != v.expectedOffset || !isRuneSliceEqual(newLine, v.expectedNewLine) {
			t.Errorf("unexpected value for case #%d expected %d %s but got %d %s", k, v.expectedOffset, runesToString(v.expectedNewLine), offset, runesToString(newLine))
		}
	}
}

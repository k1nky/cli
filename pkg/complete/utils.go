package complete

import "strconv"

func runesToString(r [][]rune) string {
	s := ""
	for _, v := range r {
		s += "," + strconv.Quote(string(v))
	}
	return s
}

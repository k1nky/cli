package parser

import (
	"strings"
)

type Parser interface {
	Parse(input string) []string
}

type QuoteParser struct{}

func (p *QuoteParser) Parse(input string) []string {
	parsed := []string{}

	input = strings.TrimSpace(input)
	chunks := strings.Fields(input)
	quoted := ""
	for _, chunk := range chunks {
		if strings.HasPrefix(chunk, "`") {
			quoted = strings.TrimLeft(chunk, "`")
		} else if strings.HasSuffix(chunk, "`") {
			quoted += " " + strings.TrimRight(chunk, "`")
			parsed = append(parsed, quoted)
			quoted = ""
		} else if len(quoted) != 0 {
			quoted += " " + chunk
		} else {
			parsed = append(parsed, chunk)
		}
	}

	return parsed
}

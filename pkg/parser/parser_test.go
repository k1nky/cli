package parser

import (
	"testing"
)

func TestParser(t *testing.T) {
	p := QuoteParser{}
	parsed := p.Parse("aa bb cc `hello world` dd")
	if len(parsed) != 5 {
		t.Errorf("Expected 5 arguments: %d", len(parsed))
	}
}

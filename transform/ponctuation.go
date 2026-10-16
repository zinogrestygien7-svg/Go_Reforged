package transform

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"Go_Reforged/lexer"
)

const (
	roleNone = iota
	roleOpen
	roleClose
	roleUnpaired
)

// Punctuation recolle la ponctuation, les guillemets et les parenthèses.
type Punctuation struct{}

func (Punctuation) Name() string { return "punctuation" }

func (Punctuation) Apply(tokens []lexer.Token) ([]lexer.Token, []lexer.Diagnostic) {
	out := make([]lexer.Token, len(tokens))
	copy(out, tokens)
	roles := quoteRoles(out)

	var diags []lexer.Diagnostic
	prevOpening, prevPunct := false, false

	for i := range out {
		tok := &out[i]
		if roles[i] == roleUnpaired {
			diags = append(diags, lexer.Diagnostic{
				Line: tok.Line, Col: tok.Col, Severity: "warning",
				Message: fmt.Sprintf("unmatched quote %q", tok.Text),
			})
		}

		separator := isNumberSeparator(out, i)
		closing := !separator && (tok.Kind == lexer.Punct || roles[i] == roleClose ||
			(tok.Kind == lexer.Paren && tok.Text == ")"))

		switch {
		case i == 0:
			// le premier token n'a rien avant lui
		case closing, prevOpening:
			tok.SpaceBefore = false
		case prevPunct:
			tok.SpaceBefore = true
		}

		prevOpening = roles[i] == roleOpen || (tok.Kind == lexer.Paren && tok.Text == "(")
		prevPunct = tok.Kind == lexer.Punct && !separator
	}
	return out, diags
}

func quoteRoles(tokens []lexer.Token) []int {
	roles := make([]int, len(tokens))
	var stack []int // index des guillemets ouverts

	for i, tok := range tokens {
		if tok.Kind != lexer.Quote {
			continue
		}
		matched := false
		for j := len(stack) - 1; j >= 0; j-- {
			if tokens[stack[j]].Text != tok.Text {
				continue
			}
			for _, k := range stack[j+1:] { // restés ouverts au-dessus : sans partenaire
				roles[k] = roleUnpaired
			}
			roles[stack[j]] = roleOpen
			roles[i] = roleClose
			stack = stack[:j]
			matched = true
			break
		}
		if !matched {
			stack = append(stack, i)
		}
	}
	for _, k := range stack {
		roles[k] = roleUnpaired
	}
	return roles
}

func isNumberSeparator(tokens []lexer.Token, i int) bool {
	tok := tokens[i]
	if tok.Kind != lexer.Punct || utf8.RuneCountInString(tok.Text) != 1 || tok.SpaceBefore {
		return false
	}
	if i == 0 || i+1 >= len(tokens) {
		return false
	}
	prev, next := tokens[i-1], tokens[i+1]
	if next.SpaceBefore || prev.Kind != lexer.Word || next.Kind != lexer.Word {
		return false
	}
	last, _ := utf8.DecodeLastRuneInString(prev.Text)
	first, _ := utf8.DecodeRuneInString(next.Text)
	return unicode.IsDigit(last) && unicode.IsDigit(first)
}

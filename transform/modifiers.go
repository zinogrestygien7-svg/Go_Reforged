package transform

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"Go_Reforged/lexer"
)

var caseFuncs = map[string]func(string) string{
	"up":  strings.ToUpper,
	"low": strings.ToLower,
	"cap": capitalize,
}

type Modifiers struct{}

func (Modifiers) Name() string { return "modifiers" }

func (Modifiers) Apply(tokens []lexer.Token) ([]lexer.Token, []lexer.Diagnostic) {
	var out []lexer.Token
	var diags []lexer.Diagnostic
	carry := false

	keep := func(tok lexer.Token) {
		tok.SpaceBefore = tok.SpaceBefore || carry
		carry = false
		out = append(out, tok)
	}
	warn := func(tok lexer.Token, msg string) {
		diags = append(diags, lexer.Diagnostic{
			Line: tok.Line, Col: tok.Col, Severity: "warning", Message: msg,
		})
	}
	for _, tok := range tokens {
		if tok.Kind != lexer.Modifier {
			keep(tok)
			continue
		}
		name, count, countOK := parseModifier(tok.Text)
		fn, supported := caseFuncs[name]

		switch {
		case !supported:
			keep(tok)
		case !countOK:
			warn(tok, fmt.Sprintf("invalid count in modifier %q", tok.Text))
			keep(tok)
		default:
			idx := lastWords(out, count)
			if len(idx) == 0 {
				warn(tok, fmt.Sprintf("no word before modifier %q", tok.Text))
				keep(tok)
				continue
			}
			if len(idx) < count {
				warn(tok, fmt.Sprintf("modifier %q asks for %d words but only %d avaible",
					tok.Text, count, len(idx)))
			}
			for _, i := range idx {
				out[i].Text = fn(out[i].Text)
				out[i].Modified = true
			}
			carry = carry || tok.SpaceBefore
		}
	}
	return out, diags
}

func parseModifier(text string) (name string, count int, countOK bool) {
	inner := strings.TrimSuffix(strings.TrimPrefix(text, "("), ")")
	rawName, rawCount, hasCount := strings.Cut(inner, ",")
	name = strings.TrimSpace(rawName)
	if !hasCount {
		return name, 1, true
	}
	n, err := strconv.Atoi(strings.TrimSpace(rawCount))
	if err != nil || n < 1 {
		return name, 0, false
	}
	return name, n, true
}

func lastWords(tokens []lexer.Token, count int) []int {
	var idx []int
	for i := len(tokens) - 1; i >= 0 && len(idx) < count; i-- {
		if tokens[i].Kind == lexer.Word {
			idx = append(idx, i)
		}
	}
	return idx
}

func capitalize(word string) string {
	runes := []rune(strings.ToLower(word))
	for i, r := range runes {
		if unicode.IsLetter(r) {
			runes[i] = unicode.ToUpper(r)
			break
		}
	}
	return string(runes)
}

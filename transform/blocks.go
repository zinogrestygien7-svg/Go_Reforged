package transform

import (
	"fmt"
	"strings"

	"Go_Reforged/lexer"
)

type Blocks struct{}

func (Blocks) name() string { return "blocks"}

type openBlock struct {
	name string
	tok lexer.Token
	words int
}

func (Blocks) Apply(tokens []lexer.Token)([]lexer.Token, []lexer.Diagnostic) {
	if diags := checkBlocks(tokens); len(diags) > 0 {
		return tokens, diags
	}
	var out []lexer.Token
	var stack []openBlock
	carry := false // espace d'une balise supprimée, à reporter sur le token suivant

	for _, tok := range tokens {
		switch tok.Kind {
		case lexer.BlockOpen:
			name, _ := blockName(tok.Text)
			stack = append(stack, openBlock{name: name, tok: tok})
			carry = carry || tok.SpaceBefore
		case lexer.BlockClose:
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
			carry = carry || tok.SpaceBefore
		default:
			tok.SpaceBefore = tok.SpaceBefore || carry
			carry = false
			if tok.Kind == lexer.Word && len(stack) > 0 {
				top := &stack[len(stack)-1] // le bloc le plus intérieur l'emporte
				tok.Text = transformInBlock(top.name, tok.Text, top.words == 0)
				tok.Modified = true
				top.words++
			}
			out = append(out, tok)
		}
	}
	return out, nil
}

func checkBlocks(tokens []lexer.Token) []lexer.Diagnostic {
	var diags []lexer.Diagnostic
	var stack []openBlock
	fail := func(tok lexer.Token, msg string) {
		diags = append(diags, lexer.Diagnostic{
			Line: tok.Line, Col: tok.Col, Severity: "error", Message: msg,
		})
	}

	for _, tok := range tokens {
		switch tok.Kind {
		case lexer.BlockOpen:
			name, _ := blockName(tok.Text)
			stack = append(stack, openBlock{name: name, tok: tok})
		case lexer.BlockClose:
			name, _ := blockName(tok.Text)
			switch {
			case len(stack) == 0:
				fail(tok, fmt.Sprintf("closing tag %q has no matching opening tag", tok.Text))
			case stack[len(stack)-1].name != name:
				fail(tok, fmt.Sprintf("closing tag %q does not match opening tag %q",
					tok.Text, stack[len(stack)-1].tok.Text))
			default:
				stack = stack[:len(stack)-1]
			}
		}
	}
	for _, b := range stack {
		fail(b.tok, fmt.Sprintf("block %q is never closed", b.tok.Text))
	}
	return diags
}

func blockName(text string) (name string, closing bool) {
	inner := strings.TrimSuffix(strings.TrimPrefix(text, "["), "]")
	if strings.HasPrefix(inner, "/") {
		return strings.TrimPrefix(inner, "/"), true
	}
	return inner, false
}

// transformInBlock applique la transformation du bloc à un mot.
// first indique si c'est le premier mot du bloc (utile pour title).
func transformInBlock(name, word string, first bool) string {
	if name == "title" {
		return titleWord(word, first)
	}
	fn, ok := modifierFuncs[name]
	if !ok {
		return word
	}
	text, _ := fn(word) // les fonctions de casse n'échouent jamais
	return text
}
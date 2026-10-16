package transform

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"Go_Reforged/lexer"
)

type Articles struct{}

func (Articles) Name() string { return "articles" }

func (Articles) Apply(tokens []lexer.Token) ([]lexer.Token, []lexer.Diagnostic) {
	out := make([]lexer.Token, len(tokens))
	copy(out, tokens)

	for i, tok := range out {
		if tok.Kind != lexer.Word {
			continue
		}
		lower := strings.ToLower(tok.Text)
		if lower != "a" && lower != "an" {
			continue
		}
		next, found := nextWord(out, i)
		if !found {
			continue
		}
		wantsAn, ok := startsWithVowelOrh(next)
		if !ok {
			continue
		}
		switch {
		case lower == "a" && wantsAn:
			out[i].Text = addN(tok.Text)
		case lower == "an" && !wantsAn:
			out[i].Text = removeN(tok.Text)
		}
	}
	return out, nil
}

func nextWord(tokens []lexer.Token, i int) (string, bool) {
	for j := i + 1; j < len(tokens); j++ {
		switch tokens[j].Kind {
		case lexer.Word:
			return tokens[j].Text, true
		case lexer.Quote, lexer.Paren:
			continue
		default:
			return "", false
		}
	}
	return "", false
}

func startsWithVowelOrh(word string) (wantsAn, ok bool) {
	r, _ := utf8.DecodeRuneInString(word)
	if !unicode.IsLetter(r) {
		return false, false
	}
	return strings.ContainsRune("aeiouhAEIOUH", r), true
}

func addN(article string) string {
	if article == "A" {
		return "An"
	}
	return article + "n"
}

func removeN(article string) string {
	if article == "An" || article == "AN" {
		return "A"
	}
	return "a"
}

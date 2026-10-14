package lexer

import (
	"strings"
	"unicode"
)

var modifierNames = map[string]bool{
	"hex": true, "bin": true, "oct": true, "to-hex": true, "to-bin": true,
	"roman": true, "to-roman": true, "up": true, "low": true, "cap": true,
	"title": true, "rev": true,
}
var blockNames = map[string]bool{"up": true, "low": true, "cap": true, "title": true}

func Lex(line string, lineNo int) []Token {
	runes := []rune(line)
	var tokens []Token
	spaceBefore := false

	add := func(kind Kind, start, end int) {
		tokens = append(tokens, Token{
			Kind: kind, Text: string(runes[start:end]),
			Line: lineNo, Col: start + 1, SpaceBefore: spaceBefore,
		})
		spaceBefore = false
	}

	for i := 0; i < len(runes); {
		r := runes[i]
		modLen := modifierLen(runes, i)
		blockKind, blockN := blockTag(runes, i)
		escLen := 0
		if r == '\\' {
			escLen = modifierLen(runes, i+1)
		}

		switch {
		case isSpace(r):
			spaceBefore = true
			i++
		case escLen > 0:
			add(Word, i+1, i+1+escLen)
			i += 1 + escLen
		case modLen > 0:
			add(Modifier, i, i+modLen)
			i += modLen
		case blockN > 0:
			add(blockKind, i, i+blockN)
			i += blockN
		case r == '(' || r == ')':
			add(Paren, i, i+1)
			i++
		case isPunct(r):
			start := i
			for i < len(runes) && isPunct(runes[i]) {
				i++
			}
			add(Punct, start, i)
		case isQuote(runes, i):
			add(Quote, i, i+1)
			i++
		default:
			start := i
			for i < len(runes) && !isWordEnd(runes, i) {
				i++
			}
			add(Word, start, i)

		}
	}
	return tokens
}

func modifierLen(runes []rune, i int) int {
	if i >= len(runes) || runes[i] != '(' {
		return 0
	}
	for j := i + 1; j < len(runes); j++ {
		switch runes[j] {
		case '(':
			return 0
		case ')':
			name, _, _ := strings.Cut(string(runes[i+1:j]), ",")
			if modifierNames[strings.TrimSpace(name)] {
				return j - i + 1
			}
			return 0

		}
	}
	return 0
}

func blockTag(runes []rune, i int) (kind Kind, n int) {
	if i >= len(runes) || runes[i] != '[' {
		return Word, 0
	}
	for j := i + 1; j < len(runes); j++ {
		if runes[j] != ']' {
			continue
		}
		name := string(runes[i+1 : j])
		closing := strings.HasPrefix(name, "/")
		if !blockNames[strings.TrimPrefix(name, "/")] {
			return Word, 0
		}
		if closing {
			return BlockClose, j - i + 1
		}
		return BlockOpen, j - i + 1
	}
	return Word, 0
}

func isWordEnd(runes []rune, i int) bool {
	r := runes[i]
	if isSpace(r) || isPunct(r) || isQuote(r) || r == '(' || r == ')' {
		return true
	}
	_, n := blockTag(runes, i)
	return n > 0
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r'
}

func isPunct(r rune) bool {
	return r == '.' || r == ',' || r == '!' || r == '?' || r == ':' || r == ';'
}

func isContraction(runes []rune, i int) bool {
	return i > 0 && i+1 < len(runes) &&
		unicode.IsLetter(runes[i-1]) && unicode.IsLetter(runes[i+1])
}

func isQuote(runes []rune, i int) bool {
	return runes[i] == '"' || (runes[i] == '\'' && !isContraction(runes, i))
}

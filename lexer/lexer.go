package lexer

import "unicode"

func Lex(line string, lineNo int) []Token {
	runes := []rune(line)
	var tokens []Token
	spaceBefore := false

	for i := 0; i < len(runes); {
		r := runes[i]

		switch {
		case isSpace(r):
			spaceBefore = true
			i++

		case isQuote(runes, i):
			tokens = append(tokens, Token{
				Kind: Quote, Text: string(r),
				Line: lineNo, Col: i + 1, SpaceBefore: spaceBefore,
			})
			spaceBefore = false
			i++
		case isPunct(r):
			start := i
			for i < len(runes) && isPunct(runes[i]) {
				i++
			}
			tokens = append(tokens, Token{
				Kind: Punct, Text: string(runes[start:i]),
				Line: lineNo, Col: start + 1, SpaceBefore: spaceBefore,
			})
			spaceBefore = false
		default:
			start := i
			for i < len(runes) && !isSpace(runes[i]) && !isPunct(runes[i]) && !isQuote(runes[i]) {
				i++
			}
			tokens = append(tokens, Token{
				Kind: Word, Text: string(runes[start:i]),
				Line: lineNo, Col: start + 1, SpaceBefore: spaceBefore,
			})
			spaceBefore = false
		}
	}
	return tokens
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

func isQuote(runes []rune, i int)bool{
	return runes[i] == '"' || (runes[i]== '\'' && !isContraction(runes, i))
}
package render

import (
	"goreforged/lexer"
	"strings"
)

func Render(tokens []lexer.Token) string {
	var b strings.Builder
	for i, tok := range tokens {
		if i > 0 && tok.SpaceBefore {
			b.WriteByte(' ')
		}
		b.WriteString(tok.Text)
	}
	return b.String()
}

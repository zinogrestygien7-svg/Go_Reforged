package transform

import (
	"Go_Reforged/lexer"
)

type Rule interface {
	Name() string
	Apply(tokens []lexer.Token) ([]lexer.Token, []lexer.Diagnostic)
}

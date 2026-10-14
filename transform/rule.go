package transform

import "goreforged/lexer"

type Rule interface {
	Name() string
	Apply(tokens []lexer.Token) ([]lexer.Token, []lexer.Diagnostic)
}

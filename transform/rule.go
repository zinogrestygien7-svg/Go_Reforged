package transform

import "Go_Reforged/lexer"

type Rule interface {
	Name() string
	Apply(tokens []lexer.Token1) ([]lexer.Token1, []lexer.Diagnostic)
}
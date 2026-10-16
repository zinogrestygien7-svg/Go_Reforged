package transform

import "Go_Reforged/lexer"

func Run(tokens []lexer.Token, rules []Rule) ([]lexer.Token, []lexer.Diagnostic) {
	var all []lexer.Diagnostic
	for _, rule := range rules {
		var diags []lexer.Diagnostic
		tokens, diags = rule.Apply(tokens)
		all = append(all, diags...)
	}
	return tokens, all
}

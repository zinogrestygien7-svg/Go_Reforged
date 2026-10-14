package lexer

import (
	"reflect"
	"testing"
)

func TestLex(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []Token
	}{
		{"mots et ponctuation", "hello, world", []Token{
			{Kind: Word, Text: "hello", Line: 1, Col: 1},
			{Kind: Punct, Text: ",", Line: 1, Col: 6},
			{Kind: Word, Text: "world", Line: 1, Col: 8, SpaceBefore: true},
		}},
		{"ligne vide", "", nil},
		{"groupe de ponctuation", "wait...", []Token{
			{Kind: Word, Text: "wait", Line: 1, Col: 1},
			{Kind: Punct, Text: "...", Line: 1, Col: 5},
		}},
		{"mot accentué", "noël ok", []Token{
			{Kind: Word, Text: "noël", Line: 1, Col: 1},
			{Kind: Word, Text: "ok", Line: 1, Col: 6, SpaceBefore: true},
		}},
		{"guillemets doubles", `He said "hi"`, []Token{
			{Kind: Word, Text: "He", Line: 1, Col: 1},
			{Kind: Word, Text: "said", Line: 1, Col: 4, SpaceBefore: true},
			{Kind: Quote, Text: `"`, Line: 1, Col: 9, SpaceBefore: true},
			{Kind: Word, Text: "hi", Line: 1, Col: 10},
			{Kind: Quote, Text: `"`, Line: 1, Col: 12},
		}},
		{"contraction", "don't", []Token{
			{Kind: Word, Text: "don't", Line: 1, Col: 1},
		}},
		{"nom propre", "O'Neil", []Token{
			{Kind: Word, Text: "O'Neil", Line: 1, Col: 1},
		}},
		{"apostrophes autour d'une contraction", "' don't '", []Token{
			{Kind: Quote, Text: "'", Line: 1, Col: 1},
			{Kind: Word, Text: "don't", Line: 1, Col: 3, SpaceBefore: true},
			{Kind: Quote, Text: "'", Line: 1, Col: 9, SpaceBefore: true},
		}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Lex(tc.in, 1)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Lex(%q)\n got: %+v\nwant: %+v", tc.in, got, tc.want)
			}
		})
	}
}

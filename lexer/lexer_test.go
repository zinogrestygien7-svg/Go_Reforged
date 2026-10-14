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
		// À toi d'ajouter : ligne vide, "wait...", "noël ok"
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

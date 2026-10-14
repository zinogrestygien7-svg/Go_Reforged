package render

import (
	"testing"
	
	"goreforged/lexer"
)

func TestRender(t *testing.T) {
	tests := []struct {
		name string
		in string
		want string
	}{
		{"texte simple", "hello, world", "hello, world"},
		{"espaces conservés", "hello ,world", "hello ,world"},
		{"espaces multiples réduits", "a   b", "a b"},
		{"ligne vide", "", ""},
		{"guillemets", `He said "hi" and left`, `He said "hi" and left`},
		{"contraction", "don't worry", "don't worry"},
		{"modificateur", "go (up, 2) now", "go (up, 2) now"},
		{"bloc", "[up]hi[/up] there", "[up]hi[/up] there"},
		{"modificateur échappé", `Type \(up) to`, "Type (up) to"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Render(lexer.Lex(tc.in, 1))
			if got != tc.want {
				t.Errorf("Render(Lex(%q))\n got: %q\nwant: %q", tc.in, got, tc.want)
			}
		})
	}
}
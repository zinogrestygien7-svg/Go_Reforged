package transform

import (
	"testing"

	"Go_Reforged/lexer"
	"Go_Reforged/render"
)

func TestArticles(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"a devant voyelle", "There it was. A amazing rock!", "There it was. An amazing rock!"},
		{"a devant h", "in a hour", "in an hour"},
		{"a devant consonne", "a big dog", "a big dog"},
		{"an devant consonne", "an banana", "a banana"},
		{"casse An", "An banana", "A banana"},
		{"casse AN", "AN banana", "A banana"},
		{"casse A", "A orange", "An orange"},
		{"chiffre", "a 3 apples", "a 3 apples"},
		{"guillemet entre les deux", `a "apple"`, `an "apple"`},
		{"ponctuation bloque", "a, apple", "a, apple"},
		{"article à la fin", "this is a", "this is a"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, _ := Articles{}.Apply(lexer.Lex(tc.in, 1))
			if got := render.Render(out); got != tc.want {
				t.Errorf("Articles(%q)\n got: %q\nwant %q", tc.in, got, tc.want)
			}
		})
	}
}

package transform

import (
	"reflect"
	"testing"

	"Go_Reforged/lexer"
	"Go_Reforged/render"
)

func TestModifiers(t *testing.T) {
	tests := []struct {
		name      string
		in        string
		want      string
		wantDiags []lexer.Diagnostic
	}{
		{"up", "Ready, set, go (up) !", "Ready, set, GO !", nil},
		{"low", "I should stop SHOUTING (low)", "I should stop shouting", nil},
		{"cap", "Welcome to the Brooklyn bridge (cap)", "Welcome to the Brooklyn Bridge", nil},
		{"up avec nombre", "This is so exciting (up, 2)", "This is SO EXCITING", nil},
		{"ponctuation ignorée dans le compte", "hello, world (up, 2)", "HELLO, WORLD", nil},
		{"lettres accentuées", "élan (cap) et ÉTÉ (low)", "Élan et été", nil},
		{"espace conservé", "hello (up)world", "HELLO world", nil},
		{"nombre trop grand", "hello (up, 5)", "HELLO", []lexer.Diagnostic{
			{Line: 1, Col: 7, Severity: "warning",
				Message: `modifier "(up, 5)" asks for 5 words but only 1 available`},
		}},
		{"nombre négatif", "go (up, -1)", "go (up, -1)", []lexer.Diagnostic{
			{Line: 1, Col: 4, Severity: "warning", Message: `invalid count in modifier "(up, -1)"`},
		}},
		{"nombre non numérique", "go (up, two)", "go (up, two)", []lexer.Diagnostic{
			{Line: 1, Col: 4, Severity: "warning", Message: `invalid count in modifier "(up, two)"`},
		}},
		{"aucun mot avant", "(up) hello", "(up) hello", []lexer.Diagnostic{
			{Line: 1, Col: 1, Severity: "warning", Message: `no word before modifier "(up)"`},
		}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, diags := Modifiers{}.Apply(lexer.Lex(tc.in, 1))
			if got := render.Render(out); got != tc.want {
				t.Errorf("Modifiers(%q)\n got: %q\nwant: %q", tc.in, got, tc.want)
			}
			if !reflect.DeepEqual(diags, tc.wantDiags) {
				t.Errorf("Modifiers(%q) diagnostics\n got: %+v\nwant: %+v", tc.in, diags, tc.wantDiags)
			}
		})
	}
}

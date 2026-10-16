package transform

import (
	"reflect"
	"testing"

	"Go_Reforged/lexer"
	"Go_Reforged/render"
)

func TestBlocks(t *testing.T) {
	tests := []struct {
		name      string
		in        string
		want      string
		wantDiags []lexer.Diagnostic
	}{
		{"up", "[up]hello world[/up]", "HELLO WORLD", nil},
		{"imbriqué", "[up]this is [low]VERY[/low] important[/up]", "THIS IS very IMPORTANT", nil},
		{"cap", "[cap]hello big world[/cap]", "Hello Big World", nil},
		{"title", "[title]the lord of the rings[/title]", "The Lord of the Rings", nil},
		{"ponctuation intacte", "[up]hello, world[/up]", "HELLO, WORLD", nil},
		{"modificateur laissé", "[low]Hello World (up)[/low]", "hello world (up)", nil},
		{"hors du bloc", "keep [up]this[/up] one", "keep THIS one", nil},
		{"jamais fermé", "[up]hello", "[up]hello", []lexer.Diagnostic{
			{Line: 1, Col: 1, Severity: "error", Message: `block "[up]" is never closed`},
		}},
		{"fermeture sans ouverture", "hello[/up]", "hello[/up]", []lexer.Diagnostic{
			{Line: 1, Col: 6, Severity: "error",
				Message: `closing tag "[/up]" has no matching opening tag`},
		}},
		{"mauvaise fermeture", "[up]a[/low]b[/up]", "[up]a[/low]b[/up]", []lexer.Diagnostic{
			{Line: 1, Col: 6, Severity: "error",
				Message: `closing tag "[/low]" does not match opening tag "[up]"`},
		}},
		{"balises croisées", "[up][low]x[/up][/low]", "[up][low]x[/up][/low]", []lexer.Diagnostic{
			{Line: 1, Col: 11, Severity: "error",
				Message: `closing tag "[/up]" does not match opening tag "[low]"`},
			{Line: 1, Col: 1, Severity: "error", Message: `block "[up]" is never closed`},
		}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, diags := Blocks{}.Apply(lexer.Lex(tc.in, 1))
			if got := render.Render(out); got != tc.want {
				t.Errorf("Blocks(%q)\n got: %q\nwant: %q", tc.in, got, tc.want)
			}
			if !reflect.DeepEqual(diags, tc.wantDiags) {
				t.Errorf("Blocks(%q) diagnostics\n got: %+v\nwant: %+v", tc.in, diags, tc.wantDiags)
			}
		})
	}
}

// Les blocs passent avant les modificateurs, qui ont donc le dernier mot.
func TestBlocksThenModifiers(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"modificateur après bloc", "[low]Hello World (up)[/low]", "hello WORLD"},
		{"exemple du sujet", "hello (up) [cap]big world[/cap]", "HELLO Big World"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, _ := Run(lexer.Lex(tc.in, 1), []Rule{Blocks{}, Modifiers{}})
			if got := render.Render(out); got != tc.want {
				t.Errorf("Run(%q)\n got: %q\nwant: %q", tc.in, got, tc.want)
			}
		})
	}
}
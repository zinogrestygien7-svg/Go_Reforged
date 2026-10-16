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
				{"hex", "1E (hex) files were added", "30 files were added", nil},
		{"hex minuscules", "ff (hex)", "255", nil},
		{"hex négatif", "-1A (hex)", "-26", nil},
		{"bin", "It has been 10 (bin) years", "It has been 2 years", nil},
		{"oct", "17 (oct)", "15", nil},
		{"to-hex", "255 (to-hex)", "FF", nil},
		{"to-bin", "5 (to-bin)", "101", nil},
		{"roman", "MMXXVI (roman)", "2026", nil},
		{"roman minuscules", "mmxxvi (roman)", "2026", nil},
		{"to-roman", "1984 (to-roman)", "MCMLXXXIV", nil},
		{"conversion avec nombre", "1E 10 (hex, 2)", "30 16", nil},
		{"chaînage", "FF (hex) (to-bin)", "11111111", nil},
		{"hex invalide", "Add ZZ (hex) now", "Add ZZ (hex) now", []lexer.Diagnostic{
			{Line: 1, Col: 8, Severity: "warning", Message: `invalid hexadecimal number "ZZ"`},
		}},
		{"bin invalide", "3 (bin)", "3 (bin)", []lexer.Diagnostic{
			{Line: 1, Col: 3, Severity: "warning", Message: `invalid binary number "3"`},
		}},
		{"roman IIII refusé", "IIII (roman)", "IIII (roman)", []lexer.Diagnostic{
			{Line: 1, Col: 6, Severity: "warning", Message: `invalid roman numeral "IIII"`},
		}},
		{"roman VX refusé", "VX (roman)", "VX (roman)", []lexer.Diagnostic{
			{Line: 1, Col: 4, Severity: "warning", Message: `invalid roman numeral "VX"`},
		}},
		{"to-roman hors intervalle", "4000 (to-roman)", "4000 (to-roman)", []lexer.Diagnostic{
			{Line: 1, Col: 6, Severity: "warning", Message: `number "4000" is out of range 1-3999`},
		}},
		{"dépassement int64", "9223372036854775808 (to-hex)", "9223372036854775808 (to-hex)", []lexer.Diagnostic{
			{Line: 1, Col: 21, Severity: "warning", Message: `invalid decimal number "9223372036854775808"`},
		}},
		{"deux mots invalides", "ZZ YY (hex, 2)", "ZZ YY (hex, 2)", []lexer.Diagnostic{
			{Line: 1, Col: 7, Severity: "warning", Message: `invalid hexadecimal number "ZZ"`},
			{Line: 1, Col: 7, Severity: "warning", Message: `invalid hexadecimal number "YY"`},
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

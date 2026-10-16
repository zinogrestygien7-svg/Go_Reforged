package transform

import (
	"reflect"
	"testing"

	"Go_Reforged/lexer"
	"Go_Reforged/render"
)

func TestPipeline(t *testing.T) {
	rules := []Rule{Blocks{}, Modifiers{}, Articles{}, Punctuation{}}

	tests := []struct {
		name      string
		in        string
		want      string
		wantDiags []lexer.Diagnostic
	}{
		{"conversions chaînées",
			"The answer is FF (hex) (to-bin) and the year is MMXXVI (roman) .",
			"The answer is 11111111 and the year is 2026.", nil},
		{"blocs",
			"[up]this is [low]VERY[/low] important[/up] , trust me .",
			"THIS IS very IMPORTANT, trust me.", nil},
		{"articles",
			"I ate an banana and a orange in a hour .",
			"I ate a banana and an orange in an hour.", nil},
		{"citations",
			`She said : " don't worry , it's ' fine ' " and left .`,
			`She said: "don't worry, it's 'fine'" and left.`, nil},
		{"avertissements",
			"Add ZZ (hex) and 3 (bin) please (up, -1) .",
			"Add ZZ (hex) and 3 (bin) please (up, -1).", []lexer.Diagnostic{
				{Line: 1, Col: 8, Severity: "warning", Message: `invalid hexadecimal number "ZZ"`},
				{Line: 1, Col: 20, Severity: "warning", Message: `invalid binary number "3"`},
				{Line: 1, Col: 33, Severity: "warning", Message: `invalid count in modifier "(up, -1)"`},
			}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, diags := Run(lexer.Lex(tc.in, 1), rules)
			if got := render.Render(out); got != tc.want {
				t.Errorf("Run(%q)\n got: %q\nwant: %q", tc.in, got, tc.want)
			}
			if !reflect.DeepEqual(diags, tc.wantDiags) {
				t.Errorf("Run(%q) diagnostics\n got: %+v\nwant: %+v", tc.in, diags, tc.wantDiags)
			}
		})
	}
}
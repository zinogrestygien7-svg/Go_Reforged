package transform

import (
	"reflect"
	"testing"

	"Go_Reforged/lexer"
	"Go_Reforged/render"
)

func TestPunctuation(t *testing.T) {
	tests := []struct {
		name      string
		in        string
		want      string
		wantDiags []lexer.Diagnostic
	}{
		{"ponctuation simple", "over there ,and then BAMM !!", "over there, and then BAMM!!", nil},
		{"groupe", "I was thinking ... You were right", "I was thinking... You were right", nil},
		{"groupe mixte", "Really !? Yes", "Really!? Yes", nil},
		{"décimal conservé", "pi is 3.14 ok", "pi is 3.14 ok", nil},
		{"heure conservée", "meet at 10:30 ok", "meet at 10:30 ok", nil},
		{"point seul avant un nombre", "wait . 5 apples", "wait. 5 apples", nil},
		{"un mot", "I am exactly how they describe me: ' awesome '",
			"I am exactly how they describe me: 'awesome'", nil},
		{"plusieurs mots", "As Elton John said: ' I am the most well-known homosexual in the world '",
			"As Elton John said: 'I am the most well-known homosexual in the world'", nil},
		{"imbriqués", `He said: " she told me ' no ' "`, `He said: "she told me 'no'"`, nil},
		{"contraction dans une citation", "' don't do it '", "'don't do it'", nil},
		{"nom propre", "I met O'Neil today", "I met O'Neil today", nil},
		{"parenthèses", "it was ( almost ) perfect", "it was (almost) perfect", nil},
		{"parenthèse puis point", "it was ( almost ) .", "it was (almost).", nil},
		{"ponctuation avant guillemet fermant", `He said " stop ! " and left`, `He said "stop!" and left`, nil},
		{"exemple du sujet", `She said : " don't worry , it's ' fine ' " and left .`,
			`She said: "don't worry, it's 'fine'" and left.`, nil},
		{"modificateur resté", "go (up, -1) !", "go (up, -1)!", nil},
		{"ligne vide", "", "", nil},
		{"guillemet double orphelin", `hello " world`, `hello " world`, []lexer.Diagnostic{
			{Line: 1, Col: 7, Severity: "warning", Message: `unmatched quote "\""`},
		}},
		{"apostrophe orpheline", "hello ' world", "hello ' world", []lexer.Diagnostic{
			{Line: 1, Col: 7, Severity: "warning", Message: `unmatched quote "'"`},
		}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, diags := Punctuation{}.Apply(lexer.Lex(tc.in, 1))
			if got := render.Render(out); got != tc.want {
				t.Errorf("Punctuation(%q)\n got: %q\nwant: %q", tc.in, got, tc.want)
			}
			if !reflect.DeepEqual(diags, tc.wantDiags) {
				t.Errorf("Punctuation(%q) diagnostics\n got: %+v\nwant: %+v", tc.in, diags, tc.wantDiags)
			}
		})
	}
}

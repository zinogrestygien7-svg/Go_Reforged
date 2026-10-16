package transform

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"Go_Reforged/lexer"
)

// wordFunc transforme un mot. L'erreur décrit pourquoi la conversion est impossible.
type wordFunc func(word string) (string, error)

// caseFunc adapte une fonction de casse (qui ne peut pas échouer) à wordFunc.
func caseFunc(f func(string) string) wordFunc {
	return func(word string) (string, error) { return f(word), nil }
}

// modifierFuncs associe chaque modificateur à sa transformation.
var modifierFuncs = map[string]wordFunc{
	"up":       caseFunc(strings.ToUpper),
	"low":      caseFunc(strings.ToLower),
	"cap":      caseFunc(capitalize),
	"hex":      fromBase(16, "hexadecimal"),
	"bin":      fromBase(2, "binary"),
	"oct":      fromBase(8, "octal"),
	"to-hex":   toBase(16),
	"to-bin":   toBase(2),
	"roman":    fromRoman,
	"to-roman": toRoman,
	"rev":      caseFunc(reverse),
}

// Modifiers applique les modificateurs postfixés : (up), (hex), (to-roman)...
type Modifiers struct{}

func (Modifiers) Name() string { return "modifiers" }

func (Modifiers) Apply(tokens []lexer.Token) ([]lexer.Token, []lexer.Diagnostic) {
	var out []lexer.Token
	var diags []lexer.Diagnostic
	carry := false // espace d'un modificateur supprimé, à reporter sur le token suivant

	keep := func(tok lexer.Token) {
		tok.SpaceBefore = tok.SpaceBefore || carry
		carry = false
		out = append(out, tok)
	}
	warn := func(tok lexer.Token, msg string) {
		diags = append(diags, lexer.Diagnostic{
			Line: tok.Line, Col: tok.Col, Severity: "warning", Message: msg,
		})
	}

	for _, tok := range tokens {
		if tok.Kind != lexer.Modifier {
			keep(tok)
			continue
		}

		name, count, countOK := parseModifier(tok.Text)
		fn, supported := modifierFuncs[name]
		isTitle := name == "title" // dépend de la position du mot : traité à part

		switch {
		case !supported && !isTitle:
			keep(tok)
		case !countOK:
			warn(tok, fmt.Sprintf("invalid count in modifier %q", tok.Text))
			keep(tok)
		default:
			idx := lastWords(out, count)
			if len(idx) == 0 {
				warn(tok, fmt.Sprintf("no word before modifier %q", tok.Text))
				keep(tok)
				continue
			}
			if len(idx) < count {
				warn(tok, fmt.Sprintf("modifier %q asks for %d words but only %d available",
					tok.Text, count, len(idx)))
			}
			var results, errs []string
			if isTitle {
				results = titleCase(wordsAt(out, idx))
			} else {
				results, errs = convertWords(out, idx, fn)
			}
			if len(errs) > 0 {
				for _, msg := range errs {
					warn(tok, msg)
				}
				keep(tok) // rien n'est modifié : le modificateur reste dans le texte
				continue
			}
			for k, i := range idx {
				out[i].Text = results[k]
				out[i].Modified = true
			}
			carry = carry || tok.SpaceBefore // le modificateur disparaît
		}
	}
	return out, diags
}

// convertWords applique fn aux mots d'index idx, sans rien modifier.
// Il renvoie les nouveaux textes et les messages d'erreur éventuels.
func convertWords(tokens []lexer.Token, idx []int, fn wordFunc) ([]string, []string) {
	results := make([]string, len(idx))
	var errs []string
	for k, i := range idx {
		text, err := fn(tokens[i].Text)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		results[k] = text
	}
	return results, errs
}

// parseModifier décompose "(up, 2)" en nom "up" et nombre 2.
// Sans nombre, le compte vaut 1. countOK vaut false si le nombre est invalide.
func parseModifier(text string) (name string, count int, countOK bool) {
	inner := strings.TrimSuffix(strings.TrimPrefix(text, "("), ")")
	rawName, rawCount, hasCount := strings.Cut(inner, ",")
	name = strings.TrimSpace(rawName)
	if !hasCount {
		return name, 1, true
	}
	n, err := strconv.Atoi(strings.TrimSpace(rawCount))
	if err != nil || n < 1 {
		return name, 0, false
	}
	return name, n, true
}

// lastWords renvoie les index des count derniers mots, dans l'ordre du texte.
func lastWords(tokens []lexer.Token, count int) []int {
	var idx []int
	for i := len(tokens) - 1; i >= 0 && len(idx) < count; i-- {
		if tokens[i].Kind == lexer.Word {
			idx = append(idx, i)
		}
	}
	for l, r := 0, len(idx)-1; l < r; l, r = l+1, r-1 {
		idx[l], idx[r] = idx[r], idx[l]
	}
	return idx
}

// capitalize met la première lettre en majuscule et le reste en minuscules.
func capitalize(word string) string {
	runes := []rune(strings.ToLower(word))
	for i, r := range runes {
		if unicode.IsLetter(r) {
			runes[i] = unicode.ToUpper(r)
			break
		}
	}
	return string(runes)
}

// fromBase convertit un nombre écrit dans la base donnée vers le décimal.
// ParseInt ignore la casse, accepte un "-" en tête et refuse ce qui dépasse int64.
func fromBase(base int, label string) wordFunc {
	return func(word string) (string, error) {
		n, err := strconv.ParseInt(word, base, 64)
		if err != nil {
			return "", fmt.Errorf("invalid %s number %q", label, word)
		}
		return strconv.FormatInt(n, 10), nil
	}
}

// toBase convertit un nombre décimal vers la base donnée (en majuscules).
func toBase(base int) wordFunc {
	return func(word string) (string, error) {
		n, err := strconv.ParseInt(word, 10, 64)
		if err != nil {
			return "", fmt.Errorf("invalid decimal number %q", word)
		}
		return strings.ToUpper(strconv.FormatInt(n, base)), nil
	}
}

var smallWords = map[string]bool{
	"a": true, "an": true, "the": true, "of": true, "and": true,
	"or": true, "in": true, "on": true, "at": true, "to": true,
}

func reverse(word string) string {
	var clusters [][]rune
	for _, r := range word {
		if unicode.IsMark(r) && len(clusters) > 0 {
			clusters[len(clusters)-1] = append(clusters[len(clusters)-1], r)
			continue
		}
		clusters = append(clusters, []rune{r})
	}
	var b strings.Builder
	for i := len(clusters) - 1; i >= 0; i-- {
		b.WriteString(string(clusters[i]))
	}
	return b.String()
}

func titleCase(words []string) []string {
	out := make([]string, len(words))
	for i, w := range words {
		lower := strings.ToLower(w)
		if i > 0 && smallWords[lower] {
			out[i] = lower
		} else {
			out[i] = capitalize(w)
		}
	}
	return out
}

func wordsAt(tokens []lexer.Token, idx []int) []string {
	words := make([]string, len(idx))
	for k, i := range idx {
		words[k] = tokens[i].Text
	}
	return words
}
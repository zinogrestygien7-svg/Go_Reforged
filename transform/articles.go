package transform

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"

	"Go_Reforged/lexer"
)

// Articles transforme "a" en "an" devant une voyelle ou un h, et inversement.
// Exceptions associe un mot (en minuscules) à l'article qu'il exige : "a" ou "an".
type Articles struct {
	Exceptions map[string]string
}

func (Articles) Name() string { return "articles" }

func (a Articles) Apply(tokens []lexer.Token) ([]lexer.Token, []lexer.Diagnostic) {
	out := make([]lexer.Token, len(tokens))
	copy(out, tokens)

	for i, tok := range out {
		if tok.Kind != lexer.Word {
			continue
		}
		lower := strings.ToLower(tok.Text)
		if lower != "a" && lower != "an" {
			continue
		}
		next, found := nextWord(out, i)
		if !found {
			continue
		}
		wantsAn, ok := a.needsAn(next)
		if !ok {
			continue
		}
		switch {
		case lower == "a" && wantsAn:
			out[i].Text = addN(tok.Text)
		case lower == "an" && !wantsAn:
			out[i].Text = removeN(tok.Text)
		}
	}
	return out, nil
}

// needsAn indique si le mot exige "an". Les exceptions passent avant la règle générale.
// Le second résultat vaut false quand aucune règle ne s'applique (un chiffre par exemple).
func (a Articles) needsAn(word string) (bool, bool) {
	if article, found := a.Exceptions[strings.ToLower(word)]; found {
		return article == "an", true
	}
	return startsWithVowelOrH(word)
}

// ParseExceptions lit un fichier d'exceptions : une ligne "mot article" par mot.
// Les lignes vides sont ignorées. Une ligne mal formée donne un diagnostic d'erreur.
func ParseExceptions(r io.Reader) (map[string]string, []lexer.Diagnostic, error) {
	exceptions := make(map[string]string)
	var diags []lexer.Diagnostic

	scanner := bufio.NewScanner(r)
	for lineNo := 1; scanner.Scan(); lineNo++ {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 0 {
			continue
		}
		article := ""
		if len(fields) == 2 {
			article = strings.ToLower(fields[1])
		}
		if len(fields) != 2 || (article != "a" && article != "an") {
			diags = append(diags, lexer.Diagnostic{
				Line: lineNo, Col: 1, Severity: "error",
				Message: fmt.Sprintf("malformed exception line %q", scanner.Text()),
			})
			continue
		}
		exceptions[strings.ToLower(fields[0])] = article
	}
	return exceptions, diags, scanner.Err()
}

// nextWord renvoie le texte du prochain mot après l'index i.
// Il saute les guillemets et les parenthèses, mais s'arrête devant le reste
// (la ponctuation, un modificateur...) : "a, apple" ne doit pas changer.
func nextWord(tokens []lexer.Token, i int) (string, bool) {
	for j := i + 1; j < len(tokens); j++ {
		switch tokens[j].Kind {
		case lexer.Word:
			return tokens[j].Text, true
		case lexer.Quote, lexer.Paren:
			continue
		default:
			return "", false
		}
	}
	return "", false
}

// startsWithVowelOrH indique si le mot commence par a, e, i, o, u ou h.
// ok vaut false si le premier caractère n'est pas une lettre (un chiffre par exemple).
func startsWithVowelOrH(word string) (wantsAn, ok bool) {
	r, _ := utf8.DecodeRuneInString(word)
	if !unicode.IsLetter(r) {
		return false, false
	}
	return strings.ContainsRune("aeiouhAEIOUH", r), true
}

func addN(article string) string { // "a" -> "an", "A" -> "An"
	if article == "A" {
		return "An"
	}
	return article + "n"
}

func removeN(article string) string { // "an" -> "a", "An" et "AN" -> "A"
	if article == "An" || article == "AN" {
		return "A"
	}
	return "a"
}

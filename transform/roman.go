package transform

import (
	"fmt"
	"strconv"
	"strings"
)

var romanSymbols = []struct {
	value  int
	symbol string
}{
	{1000, "M"}, {900, "CM"}, {500, "D"}, {400, "CD"},
	{100, "C"}, {90, "XC"}, {50, "L"}, {40, "XL"},
	{10, "X"}, {9, "IX"}, {5, "V"}, {4, "IV"}, {1, "I"},
}

var romanValues = map[rune]int{
	'I': 1, 'V': 5, 'X': 10, 'L': 50, 'C': 100, 'D': 500, 'M': 1000,
}

// intToRoman écrit n (entre 1 et 3999) en chiffres romains.
func intToRoman(n int) string {
	var b strings.Builder
	for _, s := range romanSymbols {
		for n >= s.value {
			b.WriteString(s.symbol)
			n -= s.value
		}
	}
	return b.String()
}

// romanToInt calcule la valeur d'un texte en lettres romaines, sans vérifier
// qu'il est bien formé. Un symbole plus petit que son voisin de droite se soustrait.
func romanToInt(s string) (int, bool) {
	runes := []rune(s)
	if len(runes) == 0 {
		return 0, false
	}
	total := 0
	for i, r := range runes {
		v, ok := romanValues[r]
		if !ok {
			return 0, false
		}
		if i+1 < len(runes) && v < romanValues[runes[i+1]] {
			total -= v
		} else {
			total += v
		}
	}
	return total, true
}

// fromRoman convertit un chiffre romain valide en décimal.
// Il est valide si, reconverti en romain, on retombe sur le même texte
// (cela refuse IIII et VX).
func fromRoman(word string) (string, error) {
	upper := strings.ToUpper(word)
	n, ok := romanToInt(upper)
	if !ok || n < 1 || n > 3999 || intToRoman(n) != upper {
		return "", fmt.Errorf("invalid roman numeral %q", word)
	}
	return strconv.Itoa(n), nil
}

// toRoman convertit un nombre décimal entre 1 et 3999 en chiffre romain.
func toRoman(word string) (string, error) {
	n, err := strconv.Atoi(word)
	if err != nil {
		return "", fmt.Errorf("invalid decimal number %q", word)
	}
	if n < 1 || n > 3999 {
		return "", fmt.Errorf("number %q is out of range 1-3999", word)
	}
	return intToRoman(n), nil
}
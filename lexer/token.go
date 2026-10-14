package lexer

type Kind int

const (
	Mot Kind = iota
	Punct
	Quote
	Modifier
	BlockOuvert
	BlockFerme
)

type Token1 struct {
	Kind    Kind
	Text    string
	Ligne   int
	Colonne int
}
type Diagnostic struct {
	Ligne   int
	Colonne int
	Erreur  string
	Message string
}

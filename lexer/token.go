package lexer 

type Token0 int 

const (
	Mot 		Token0 = iota
	Punct
	Quote
	Modifier
	BlockOuvert
	BlockFerme
)

type Token1 struct {
	Token0 Token0
	Text string
	Ligne int
	Colonne int
}
 type Diagnostic struct {
	Ligne int
	Colonne int
	Erreur string
	Message string
 }
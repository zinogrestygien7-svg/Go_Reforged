package lexer

// Kind indique la nature d'un token.
type Kind int

const (
	Word       Kind = iota // un mot ou un nombre
	Punct                  // un groupe de . , ! ? : ;
	Quote                  // ' ou "
	Paren                  // ( ou ) qui n'appartient pas à un modificateur
	Modifier               // (up), (hex, 3), ...
	BlockOpen              // [up], [cap], ...
	BlockClose             // [/up], [/cap], ...
)

// Token est un morceau de texte avec sa position (à partir de 1, en runes).
type Token struct {
	Kind        Kind
	Text        string
	Line        int
	Col         int
	SpaceBefore bool // vrai si un espace précédait ce token dans le texte d'origine
}

// Diagnostic décrit un avertissement ou une erreur à une position précise.
type Diagnostic struct {
	Line     int
	Col      int
	Severity string // "warning" ou "error"
	Message  string
}

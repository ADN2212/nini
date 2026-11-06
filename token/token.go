package token

// Los tokens son las unidades sematicas del lenguaje
type Token struct {
	Type    TokenType
	Lexeme  string //un lexema es el string que represanta el token
	Literal any    //el valor real que tendra el token dentro del programa
	Line    int    //la linea en la que se encuentra
}

type TokenType int

const (
	// Single-character tokens.
	LeftParen  = iota // (
	RightParen        // )
	LeftBrace         // {
	RightBrace        // }
	Comma             // ,
	Dot               // .
	Minus             // -
	Plus              // +
	Semicolon         // ;
	Slash             // /
	Star              // *

	// One or two character tokens.
	Bang         // !
	BangEqual    // !=
	Equal        // =
	EqualEqual   // ==
	Greater      // >
	GreaterEqual // >=
	Less         // <
	LessEqual    // <=

	// Literals.
	Identifier // e.g. myVar, _count
	String     // e.g. "hello"
	Number     // e.g. 123, 3.14

	// Keywords.
	And   // and
	Else  // else
	False // false
	//Fun    // fun
	//For    // for
	If    // if
	Nil   // nil
	Or    // or
	Print // print
	//Return // return
	//Super  // super
	//This   // this
	//True   // true
	Var   // var
	While // while
	EOF   // end of input (no lexeme)
)

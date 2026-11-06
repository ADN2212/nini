package lexer

import (
	"github.com/ADN2212/nini/token"
	"strconv"
)

type lexer struct {
	src     []rune //este es el texto del codigo fuente, recuerda que una string es un slice de runes.
	tokens  []token.Token
	err     error
	start   int
	current int
	line    int
}

func ScanTokens(sourceCode string) ([]token.Token, error) {
	lexer := lexer{src: []rune(sourceCode)}

	for !lexer.isAtEnd() {
		lexer.start = lexer.current
		lexer.scanNextToken() //se escanea un token a la vez
		if lexer.err != nil {
			return []token.Token{}, lexer.err //early return para mantener las cosas simples
		}
	}
	return lexer.tokens, nil
}

func (lexer *lexer) isAtEnd() bool {
	return lexer.current >= len(lexer.src)
}

func (lexer *lexer) scanNextToken() {

	char := lexer.advance()

	switch char {
	case '(':
		lexer.addToken(token.LeftParen, nil) //algunos tokens no tienen valor literal
	case ')':
		lexer.addToken(token.RightParen, nil)
	case '{':
		lexer.addToken(token.LeftBrace, nil)
	case '}':
		lexer.addToken(token.RightBrace, nil)
	case ',':
		lexer.addToken(token.Comma, nil)
	case '.':
		lexer.addToken(token.Dot, nil)
	case '-':
		lexer.addToken(token.Minus, nil)
	case '+':
		lexer.addToken(token.Plus, nil)
	case ';':
		lexer.addToken(token.Semicolon, nil)
	case '*':
		lexer.addToken(token.Star, nil)

	//en nini usare "no" para hacer negacioens
	case '!':
		if lexer.match('=') {
			lexer.addToken(token.BangEqual, nil)
		} else {
			lexer.err = ErrUnexpectedCharacter //el bang solo no es un caracter permitido.
		}
	case '=':
		if lexer.match('=') {
			lexer.addToken(token.EqualEqual, nil)
		} else {
			lexer.addToken(token.Equal, nil)
		}

	case '<':
		if lexer.match('=') {
			lexer.addToken(token.LessEqual, nil)
		} else {
			lexer.addToken(token.Less, nil)
		}

	case '>':
		if lexer.match('=') {
			lexer.addToken(token.GreaterEqual, nil)
		} else {
			lexer.addToken(token.Greater, nil)
		}
	case '/':
		// If you encounter two slashes in a row, consume a comment until the end of the line
		//observeces como estos dos metodos no agregan ningun token
		//esto es porque al interprete no le importan los comentarios, osea, estos no participan en la ejecuion de codigo.
		if lexer.match('/') {
			lexer.consumeLineComment()
			break
		}
		//no multiline comments por ahora.
		lexer.addToken(token.Slash, nil)
	case ' ', '\r', '\t':
		//Ignorar todos los espacios en blanco
	case '\n':
		lexer.line++ //si encuetra un salto de linea se limita a sumar uno a la linea

	case '"':
		lexer.scanString()
	default:
		if isDigit(char) {
			lexer.scanNumber()
		} else if isAlpha(char) {
			lexer.scanIdentifier()
		} else {
			lexer.err = ErrUnexpectedCharacter
		}

	}
}

var keywords map[string]token.TokenType = map[string]token.TokenType{
	"y":         token.And,
	"sino":      token.Else,
	"falso":     token.False,
	"si":        token.If,
	"nulo":      token.Nil,
	"o":         token.Or,
	"ver":       token.Print,
	"verdadero": token.True,
	"variable":  token.Var,
	"mientras":  token.While,
	"no":        token.Not,
}

// Los indentificadores son los tokens que sirven para
// nombrar variables, ej: num, hola10
func (lexer *lexer) scanIdentifier() {

	for isAlphaNumeric(lexer.peek()) {
		lexer.advance()
	}

	//si es una palabra clave agregar ese token como tal
	text := string(lexer.src[lexer.start:lexer.current])
	if t, ok := keywords[text]; ok {
		lexer.addToken(t, nil) //las palabras claves no tienen valor asignado
		return
	}

	//de lo contrario lo agregamos como identifier
	//look estos no tienen valor definido
	lexer.addToken(token.Identifier, nil) //los identificadores si, pero no en esta estapa
}

func isAlphaNumeric(c rune) bool {
	return isAlpha(c) || isDigit(c)
}

// como peek pero ve el que le sigue al siguiente
func (lexer *lexer) peekNext() rune {
	next := lexer.current + 1
	//si se cumple esto es porque que esta al final del file (o de la linea)
	if next >= len(lexer.src) {
		return 0
	}
	return lexer.src[next]
}

func (lexer *lexer) scanNumber() {

	//consumir todos los digitos
	for isDigit(lexer.peek()) {
		lexer.advance()
	}

	//si encuentra un punto y el caracter siguiente es un digito
	//entonces estamos frete a un rual, ej: 5.55
	if lexer.peek() == '.' && isDigit(lexer.peekNext()) {
		// consume the dot
		lexer.advance()
		//otra vez consumir todos los digitos encontrados.
		for isDigit(lexer.peek()) {
			lexer.advance()
		}
	}

	//finalmente toma el string (ej: "4.65") y lo parsea a flotante
	text := string(lexer.src[lexer.start:lexer.current])
	num, err := strconv.ParseFloat(text, 64)
	if err != nil {
		lexer.err = ErrInvalidNumberLiteral
		return
	}
	//al final agrega un token de tipo num
	lexer.addToken(token.Number, num)
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c rune) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_'
}

// esta funcion toma el caracter al que se apunta actualmente
// luego da un paso hacia adelante
// current es el siguiente caracter a consumir, no el ultimo que se consumio
func (lexer *lexer) advance() rune {
	result := lexer.src[lexer.current]
	lexer.current = lexer.current + 1
	return result
}

func (lexer *lexer) addToken(t token.TokenType, literal any) {
	//Aqui se esta extrayendo del sub string que nombra el token, ej: "variable"
	lexeme := string(lexer.src[lexer.start:lexer.current])
	lexer.tokens = append(lexer.tokens, token.Token{
		Type:    t,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    lexer.line,
	})
}

// Esta funcion sirve para detectar tokens compuestos como !=, >=, == etc
// al hallar por ejemplo un !, antes de concluir que es un not, debemos comprobar que no sea un !=
// lo mismo con > para >= y < para <=
func (lexer *lexer) match(expected rune) bool {
	if lexer.isAtEnd() {
		return false
	}
	if lexer.src[lexer.current] != expected {
		return false
	}
	lexer.current = lexer.current + 1
	return true
}

func (lexer *lexer) consumeLineComment() {
	//esto es como un while,
	//mientras el token siguiente no sea un salto de linea (o el final del file) se seguira consumiento el siguiente carracter
	for lexer.peek() != '\n' && !lexer.isAtEnd() {
		lexer.advance()
	}
}

// este metodo "chekea" el caracter siguiente pero no lo consume
func (lexer *lexer) peek() rune {
	if lexer.isAtEnd() {
		return 0
	}
	return lexer.src[lexer.current]
}

func (lexer *lexer) scanString() {
	//agregar caracteres al string hasta que se encuente el cierre "
	//o que se llegue al final del file
	for lexer.peek() != '"' && !lexer.isAtEnd() {
		//esto es para soportar saltos de linea dentro de los strings (no se que tan buena idea sea)
		if lexer.peek() == '\n' {
			lexer.line++
		}
		lexer.advance()
	}

	if lexer.isAtEnd() {
		lexer.err = ErrUnterminatedString
		return
	}

	lexer.advance()
	//fmt.Println("Despues del segundo advance ...")
	//agarra el subslice y lo transforma en string.
	value := string(lexer.src[lexer.start+1 : lexer.current-1])
	lexer.addToken(token.String, value) //genera el token corrspondente
}

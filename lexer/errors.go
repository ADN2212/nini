package lexer

import "errors"

var ErrUnterminatedString = errors.New("cadena de texto sin terminar.")
var ErrInvalidNumberLiteral = errors.New("literal numerico invilido.")
var ErrUnexpectedCharacter = errors.New("Caracter desconosido")

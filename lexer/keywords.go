package lexer

import "github.com/ADN2212/nini/token"

var keywords map[string]token.TokenType = map[string]token.TokenType{
	"y":         token.And,
	"sino":      token.Else,
	"falso":     token.False,
	"si":        token.If,
	"nulo":      token.Nil,
	"o":         token.Or,
	"ver":       token.Print,
	"verdadero": token.True,
	"sea":       token.Var,
	"mientras":  token.While,
	"no":        token.Not,
}

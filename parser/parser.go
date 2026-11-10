package parser

import (
	"errors"
	"github.com/ADN2212/nini/ast"
	"github.com/ADN2212/nini/token"
)

var UnknowPrymaryTypeErr = errors.New("Tipo primario desconocido")

type parser struct {
	tokens  []token.Token
	current int //este es el indice del token actual
	err     error
}

func Parse(tokens []token.Token) ([]ast.Stmt, error) {
	statements := []ast.Stmt{}
	parser := parser{tokens: tokens}

	for !parser.isAtEnd() {
		newStmt := parser.declaration()
		if parser.err != nil {
			return []ast.Stmt{}, parser.err
		}
		statements = append(statements, newStmt)
	}
	return statements, nil
}

// este es el top del grammar, todo es una declaracion
// sin embargo desde un punto de vista de implementacion todo es un statement
func (p *parser) declaration() ast.Stmt {
	return p.statement()
}

func (p *parser) statement() ast.Stmt {

	if p.match(token.Print) {
		expr := p.expression()
		p.consume(token.Semicolon, "Se espera ';' despues del valor")
		return ast.NewPrintStmt(&expr)
	}

	expr := p.expression()
	stmt := ast.NewExprStmt(&expr)
	p.consume(token.Semicolon, "Se espera ';' despues dela expresion")
	return stmt
}

func (p *parser) expression() ast.Expr {
	return p.equality()
}

func (p *parser) equality() ast.Expr {
	expr := p.comparison()
	newExpr := ast.Expr{} //this is a sort of hack, I am missing the "why" here.
	for p.match(token.BangEqual, token.EqualEqual) {
		operator := p.previous()
		r := p.comparison()
		newExpr = ast.NewBynaryExpr(operator, &expr, &r, ast.EqualityExpr)
	}
	if newExpr.IsEmpty() {
		return expr
	}
	return newExpr
}

func (p *parser) comparison() ast.Expr {
	expr := p.term()
	newExpr := ast.Expr{}
	for p.match(token.Greater, token.GreaterEqual, token.Less, token.LessEqual) {
		operator := p.previous()
		r := p.term()
		newExpr = ast.NewBynaryExpr(operator, &expr, &r, ast.ComparisonExpr)
	}
	if newExpr.IsEmpty() {
		return expr
	}
	return newExpr
}

func (p *parser) term() ast.Expr {
	expr := p.factor()
	newExpr := ast.Expr{}
	for p.match(token.Minus, token.Plus) {
		operator := p.previous()
		r := p.factor()
		newExpr = ast.NewBynaryExpr(operator, &expr, &r, ast.TermExpr)
	}
	if newExpr.IsEmpty() {
		return expr
	}
	return newExpr
}

func (p *parser) factor() ast.Expr {
	expr := p.unary()
	newExpr := ast.Expr{}
	for p.match(token.Slash, token.Star) {
		operator := p.previous()
		right := p.unary()
		newExpr = ast.NewBynaryExpr(operator, &expr, &right, ast.FactorExpr)
	}
	if newExpr.IsEmpty() {
		return expr
	}
	return newExpr
}

func (p *parser) unary() ast.Expr {
	if p.match(token.Not, token.Minus) {
		operator := p.previous()
		right := p.unary() //esto es lo que determina la asociatividad a la derecha; !!true = !(!true)
		return ast.NewUnaryExpr(operator, &right)
	}
	return p.primary()
}

func (p *parser) primary() ast.Expr {
	expr := ast.Expr{Type: ast.PrimaryExpr}
	if p.match(token.False) {
		expr.Val = false
		return expr
	}

	if p.match(token.True) {
		expr.Val = true
		return expr
	}

	if p.match(token.Nil) {
		expr.Val = nil
		return expr
	}

	if p.match(token.Number, token.String) {
		expr.Val = p.previous().Literal
		return expr
	}

	if p.match(token.LeftParen) {
		innerExpr := p.expression()
		expr = ast.NewGroupExpr(&innerExpr)
		p.consume(token.RightParen, "se esepra ')' al final de la exprecion grupal")
		return expr
	}

	p.err = UnknowPrymaryTypeErr
	return expr
}

func (p *parser) consume(t token.TokenType, message string) token.Token {
	if p.check(t) {
		return p.advance()
	}
	p.err = errors.New(message)
	return token.Token{}
}

func (p *parser) match(types ...token.TokenType) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *parser) check(tokenType token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *parser) peek() token.Token {
	return p.tokens[p.current]
}

func (p *parser) previous() token.Token {
	return p.tokens[p.current-1]
}

func (p *parser) peekNext() token.Token {
	return p.tokens[p.current+1]
}

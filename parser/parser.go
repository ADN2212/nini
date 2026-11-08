package parser

import (
	"errors"

	"github.com/ADN2212/nini/ast"
	"github.com/ADN2212/nini/token"
)

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

// este es el top del grammar
func (p *parser) declaration() ast.Stmt {
	return p.statement()
}

func (p *parser) statement() ast.Stmt {
	expr := p.expression()
	stmt := ast.NewExprStatement(&expr)
	return stmt
}

func (p *parser) expression() ast.Expr {
	return p.equality()
}

func (p *parser) equality() ast.Expr {
	expr := p.comparison()
	for p.match(token.BangEqual, token.EqualEqual) {
		operator := p.previous()
		r := p.comparison()
		expr = ast.NewBynaryExpr(operator, &expr, &r, ast.EqualityExpr)
	}
	return expr
}

func (p *parser) comparison() ast.Expr {
	expr := p.term()
	for p.match(token.Greater, token.GreaterEqual, token.Less, token.LessEqual) {
		operator := p.previous()
		r := p.term()
		expr = ast.NewBynaryExpr(operator, &expr, &r, ast.ComparisonExpr)
	}
	return expr
}

func (p *parser) term() ast.Expr {
	expr := p.factor()
	for p.match(token.Minus, token.Plus) {
		operator := p.previous()
		r := p.factor()
		expr = ast.NewBynaryExpr(operator, &expr, &r, ast.TermExpr)
	}
	return expr
}

func (p *parser) factor() ast.Expr {
	expr := p.unary()		
	for p.match(token.Slash, token.Star) {
		operator := p.previous()
		right := p.unary()
		//expr = ast.NewFactorExpr(operator, &expr, &right)
		expr = ast.NewBynaryExpr(operator, &expr, &right, ast.FactorExpr)
	}
	return expr
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
		expr.Val = p.previous()
		return expr
	}

	//al hayar un '(' se empieza a calcular una nueva expression
	if p.match(token.LeftParen) {
		expr = p.expression()
		p.consume(token.RightParen, "se esepra ')' al final de la exprecion grupal")
		return expr
	}

	p.err = errors.New("Tipo primario desconocido")
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

package parser

import (
	"errors"
	"fmt"

	"github.com/ADN2212/nini/ast"
	"github.com/ADN2212/nini/token"
)

type parser struct {
	tokens  []token.Token
	current int //este es el indice del token actual
	err     error
}

func Parse(tokens []token.Token) ([]ast.Stmt, error) {
	fmt.Println("Parsing")
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
	fmt.Println("In declaration")
	return p.statement()
}

func (p *parser) statement() ast.Stmt {
	fmt.Println("In statement")
	expr := p.expression()
	stmt := ast.NewExprStatement(&expr)
	return stmt
}

func (p *parser) expression() ast.Expr {
	fmt.Println("In expression")
	return p.equality()
}

func (p *parser) equality() ast.Expr {
	fmt.Println("In equality")
	expr := p.comparison()
	for p.match(token.BangEqual, token.EqualEqual) {
		operator := p.previous()
		r := p.comparison()
		expr = ast.NewBynaryExpr(operator, &expr, &r, ast.EqualityExpr)
	}
	return expr
}

func (p *parser) comparison() ast.Expr {
	fmt.Println("In comparison")
	expr := p.term()
	for p.match(token.Greater, token.GreaterEqual, token.Less, token.LessEqual) {
		operator := p.previous()
		r := p.term()
		expr = ast.NewBynaryExpr(operator, &expr, &r, ast.ComparisonExpr)
	}
	return expr
}

func (p *parser) term() ast.Expr {
	fmt.Println("In term")
	expr := p.factor()
	for p.match(token.Minus, token.Plus) {
		operator := p.previous()
		r := p.factor()
		expr = ast.NewBynaryExpr(operator, &expr, &r, ast.TermExpr)
	}
	return expr
}

func (p *parser) factor() ast.Expr {
	fmt.Println("In factor")
	expr := p.unary()		
	for p.match(token.Slash, token.Star) {
		operator := p.previous()
		right := p.unary()
		expr = ast.NewBynaryExpr(operator, &expr, &right, ast.FactorExpr)
	}
	return expr
}

func (p *parser) unary() ast.Expr {
	fmt.Println("In unary")
	if p.match(token.Not, token.Minus) {
		fmt.Println("Dentro")
		operator := p.previous()
		right := p.unary() //esto es lo que determina la asociatividad a la derecha; !!true = !(!true)
		return ast.NewUnaryExpr(operator, &right)
	}
	return p.primary()
}

func (p *parser) primary() ast.Expr {
	fmt.Println("In prymary")
	expr := ast.Expr{Type: ast.PrimaryExpr}

	if p.match(token.False) {
		fmt.Println("In false")
		expr.Val = false
		return expr
	}

	if p.match(token.True) {
		fmt.Println("In true")
		expr.Val = true
		return expr
	}

	if p.match(token.Nil) {
		fmt.Println("Un nil")
		expr.Val = nil
		return expr
	}

	fmt.Println("Before")
	if p.match(token.Number, token.String) {
		fmt.Println("Number or string")
		fmt.Println(p.current)
		expr.Val = p.previous()
		return expr
	}

	fmt.Println("llega aqui ???")
	//al hayar un '(' se empieza a calcular una nueva expression
	if p.match(token.LeftParen) {
		innerExpr := p.expression()
		expr = ast.NewGroupExpr(&innerExpr) 
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
	fmt.Println("mathc")
	// if p.isAtEnd() {
	// 	return false
	// }
	for i, tokenType := range types {
		fmt.Println(i)
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *parser) check(tokenType token.TokenType) bool {
	fmt.Println("check")
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *parser) advance() token.Token {
	fmt.Println("advance")
	if !p.isAtEnd() {
		fmt.Println("at end")
		p.current++
	}
	return p.previous()
}

func (p *parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *parser) peek() token.Token {
	fmt.Println("peek")
	return p.tokens[p.current]
}

func (p *parser) previous() token.Token {
	fmt.Println("previous")
	return p.tokens[p.current-1]
}

func (p *parser) peekNext() token.Token {
	return p.tokens[p.current+1]
}

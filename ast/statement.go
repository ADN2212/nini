package ast

import (
	"github.com/ADN2212/nini/token"
)

type Stmt struct {
	Type       StmtType
	Name       *token.Token //en caso de que sea la declaracion de una variable
	Expr       *Expr
	Condition  *Expr //En caso de que sea un if o un while
	ThenBranch *Stmt
	ElseBranch *Stmt
}

func NewExprStmt(expr *Expr) Stmt {
	return Stmt{Type: ExprStmt, Expr: expr}
}

func NewPrintStmt(expr *Expr) Stmt {
	return Stmt{Type: PrintStmt, Expr: expr}
}

func NewVarStm(name token.Token, expr *Expr) Stmt {
	return Stmt{Name: &name, Expr: expr, Type: VarDeclarationStmt}
}

type Expr struct {
	Val      any
	Type     ExprType
	Left     *Expr
	Right    *Expr
	Operator *token.Token
}

// Comprueba que una expression inicializda este vacia, es decir que no tenga valores asignados a ninguno de los campos.
func (e *Expr) IsEmpty() bool {
	return e.Val == nil && e.Type == 0 && e.Left == nil && e.Right == nil && e.Operator == nil
}

func NewUnaryExpr(operator token.Token, rightExpr *Expr) Expr {
	return Expr{
		Type:     UnaryExpr,
		Operator: &operator,
		Right:    rightExpr,
	}
}

func NewBynaryExpr(
	operator token.Token,
	leftExpr *Expr,
	rightExpr *Expr,
	_type ExprType) Expr {
	return Expr{
		Type:     _type,
		Operator: &operator,
		Left:     leftExpr,
		Right:    rightExpr,
	}
}

func NewGroupExpr(expr *Expr) Expr {
	return Expr{
		Type:  GroupExpr,
		Right: expr, //podria ser cualquiera de los dos
	}
}

type StmtType int

const (
	VarDeclarationStmt = iota + 1
	ExprStmt
	PrintStmt
	IfStmt
	WhileStmt
	BlockStmt
)

type ExprType int

// este orden va de la mas compleja a la mas simple
const (
	AssingmentExpr = iota + 1
	LogicalExpr
	EqualityExpr   //
	ComparisonExpr //
	TermExpr       //
	FactorExpr     //
	UnaryExpr      //
	PrimaryExpr    //estos son los literals
	GroupExpr      //
)

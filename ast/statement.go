package ast

import (
	"github.com/ADN2212/nini/token"
)

type Stmt struct {
	Type StmtType
	Name *token.Token//en caso de que sea la declaracion de una variable
	Expr *Expr
	Condition *Expr//En caso de que sea un if o un while
	ThenBranch *Stmt
	ElseBranch *Stmt
}

type Expr struct {
	Val any
	Type ExprType
	Left *Expr
	Right *Expr
	Operator *token.Token
}

type StmtType int

const (
	VarStmt = iota + 1
	ExprStmt
	PrintStmt
	IfStmt
	WhileStmt
	BlockStmt
)

type ExprType int

//este orden va de la mas compleja a la mas simple
const (
	AssingmentExpr = iota + 1
	LogicalExpr
	EqualityExpr
	ComparisonExpr
	TermExpr
	FactorExpr
	UnaryExpr
	PrimaryExpr
)

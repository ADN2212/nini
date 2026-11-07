package ast

import (
	"github.com/ADN2212/nini/token"
)



type Stmt struct {
	Type int
}





type Expr struct {
	Val any
	Type ExprType
	Left *Expr
	Right *Expr
	Operator *token.Token
}

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

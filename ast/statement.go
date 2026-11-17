package ast

import (
	"github.com/ADN2212/nini/token"
)

type StmtType int

const (
	VarDeclarationStmt = iota + 1
	VarStmt
	ExprStmt
	PrintStmt
	IfStmt
	WhileStmt
	BlockStmt
)

type Stmt struct {
	Type       StmtType
	Name       *token.Token //en caso de que sea la declaracion de una variable
	Expr       *Expr       
	Condition  *Expr        //En caso de que sea un if o un while
	ThenBranch *Stmt//Este tambien puede ser el body de un while loop
	ElseBranch *Stmt
	Statements []Stmt //esto es para los blocks
}


func NewWhileStmt(condition *Expr, body *Stmt) Stmt {
	return Stmt {
		Type: WhileStmt,
		Condition: condition,
		ThenBranch: body,//esta es la expression que se ejecuta en cada ejecucion del body.
	}
}

func NewIfStmt(condition *Expr, thenBranch *Stmt, elseBranch *Stmt) Stmt {
	return Stmt{
		Type:       IfStmt,
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

func NewBlockStmt(statements []Stmt) Stmt {
	return Stmt{
		Type:       BlockStmt,
		Statements: statements,
	}
}

func NewExprStmt(expr *Expr) Stmt {
	return Stmt{Type: ExprStmt, Expr: expr}
}

func NewPrintStmt(expr *Expr) Stmt {
	return Stmt{Type: PrintStmt, Expr: expr}
}

func NewVarDeclarationStm(name token.Token, expr *Expr) Stmt {
	return Stmt{Name: &name, Expr: expr, Type: VarDeclarationStmt}
}

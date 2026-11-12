package interpreter

import (
	"errors"
	"fmt"
	"github.com/ADN2212/nini/ast"
	"github.com/ADN2212/nini/enviroment"
	"github.com/ADN2212/nini/token"
)

// is this the right place ???
var env enviroment.Env = enviroment.NewEnv()

// esto es para fines de representacion, osea, lo que el user del lenguage ve.
var goToNiniRprMap map[any]string = map[any]string{
	false: "falso",
	true:  "verdadero",
}

// mover esto a un file a parte
var NotImplementedErr error = errors.New("Esta vaina no ha sido implementada")
var UnknownStmtTypeErr error = errors.New("Tipo de statement desconocido")
var UnknowExprTypeErr error = errors.New("Tipo de expression desconocido")
var ShouldBeNumberErr error = errors.New("Esto deberia ser un numero")
var UnknowUnaryOperator error = errors.New("Operador unario desconocido")
var UnknowBynaryOperator error = errors.New("Operador binario desconocido")
var UnknowBynaryOperation error = errors.New("Operacion binaria desconocida")

func Interpret(stmts []ast.Stmt) (any, error) {
	for _, stmt := range stmts {
		switch stmt.Type {
		//el ordern importa aqui ???
		case ast.VarDeclarationStmt:
			err := declareVar(stmt)
			if err != nil {
				return nil, err
			}
		case ast.PrintStmt:
			err := printExpr(stmt.Expr)
			if err != nil {
				return nil, err
			}
		case ast.IfStmt:
			return nil, NotImplementedErr
		case ast.WhileStmt:
			return nil, NotImplementedErr
		case ast.BlockStmt:
			return nil, NotImplementedErr
		case ast.ExprStmt:
			val, err := computeExpr(stmt.Expr)
			if err != nil {
				return nil, err
			}
			v, ok := goToNiniRprMap[val]
			if ok {
				val = v
			}
			fmt.Printf("expresion ejecutada => %v \n", val)
		default:
			return nil, UnknownStmtTypeErr
		}
	}
	return nil, nil
}

func declareVar(stmt ast.Stmt) error {
	val, err := computeExpr(stmt.Expr)
	if err != nil {
		return err
	}
	env.Define(stmt.Name.Lexeme, val)
	fmt.Printf("variable '%v' declarada con valor '%v' \n", stmt.Name.Lexeme, val)
	return nil
}

func printExpr(expr *ast.Expr) error {
	val, err := computeExpr(expr)
	if err != nil {
		return err
	}
	v, ok := goToNiniRprMap[val]
	if ok {
		val = v
	}
	fmt.Printf("mira => %v \n", val)
	return nil
}

func computeExpr(expr *ast.Expr) (any, error) {
	switch expr.Type {
	case ast.AssingmentExpr:
		return nil, NotImplementedErr
	//Estas 4 caen dentro de la cantegoria de operaciones binarias
	case ast.LogicalExpr:
		return nil, NotImplementedErr
	case ast.EqualityExpr:
		return computeBynary(expr)
	case ast.ComparisonExpr:
		return computeBynary(expr)
	case ast.TermExpr:
		return computeBynary(expr)
	case ast.FactorExpr:
		return computeBynary(expr)
	//--------------------------------------------------------------
	case ast.UnaryExpr:
		return computeUnary(expr)
	case ast.PrimaryExpr:
		return expr.Val, nil
	case ast.GroupExpr:
		return computeExpr(expr.Right)

	case ast.VarExpr:
		return computeVarExpr(expr)
	default:
		return nil, UnknowExprTypeErr
	}
}

func computeVarExpr(expr *ast.Expr) (any, error) {
	return env.Get(*expr.Name)
}

func computeUnary(expr *ast.Expr) (any, error) {
	r, err := computeExpr(expr.Right)
	if err != nil {
		return nil, err
	}
	switch expr.Operator.Type {
	case token.Minus:
		num, ok := r.(float64)
		if !ok {
			return nil, ShouldBeNumberErr
		}
		return -num, nil
	case token.Not: //esto es el no
		return !isTruthy(r), nil
	default:
		return nil, UnknowUnaryOperator
	}
}

func computeBynary(expr *ast.Expr) (any, error) {

	l, err := computeExpr(expr.Left)
	if err != nil {
		return 0, err
	}
	r, err := computeExpr(expr.Right)
	if err != nil {
		return 0, err
	}

	leftVal, okL := l.(float64)
	rightVal, okR := r.(float64)

	if okL && okR {
		switch expr.Operator.Type {
		case token.Plus:
			return leftVal + rightVal, nil
		case token.Minus:
			return leftVal - rightVal, nil
		case token.Slash:
			return leftVal / rightVal, nil
		case token.Star:
			return leftVal * rightVal, nil
		case token.Greater:
			return leftVal > rightVal, nil
		case token.GreaterEqual:
			return leftVal >= rightVal, nil
		case token.Less:
			return leftVal < rightVal, nil
		case token.LessEqual:
			return leftVal <= rightVal, nil
		case token.BangEqual: //!=
			return leftVal != rightVal, nil
		case token.EqualEqual:
			return leftVal == rightVal, nil
		default:
			return nil, UnknowBynaryOperator
		}
	}

	leftStr, okL := l.(string)
	rightStr, okR := r.(string)
	if okL && okR {
		switch expr.Operator.Type {
		case token.Plus:
			return leftStr + rightStr, nil
		case token.BangEqual: //!=
			return leftVal != rightVal, nil
		case token.EqualEqual:
			return leftVal == rightVal, nil
		default:
			return 0, UnknowBynaryOperator
		}
	}

	leftBool, okL := l.(bool)
	rightBool, okR := r.(bool)
	if okL && okR {
		switch expr.Operator.Type {
		case token.BangEqual:
			return leftBool != rightBool, nil
		case token.EqualEqual:
			return leftBool == rightBool, nil
		default:
			return nil, UnknowBynaryOperator
		}
	}
	return nil, UnknowBynaryOperation
}

// this is a helper
// ojo, la truthyness so podria definir de forma diferente.
func isTruthy(val any) bool {
	if val == nil {
		return false
	}
	b, ok := val.(bool)
	if ok {
		return b
	}
	return true
}

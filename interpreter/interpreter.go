package interpreter

import (
	"fmt"
	"github.com/ADN2212/nini/ast"
	"github.com/ADN2212/nini/enviroment"
	"github.com/ADN2212/nini/token"
)

// esto es para fines de representacion, osea, lo que el user del lenguage ve.
var goToNiniRprMap map[any]string = map[any]string{
	false: "falso",
	true:  "verdadero",
}

func execute(stmt ast.Stmt) (any, error) {
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
		err := computeIfStmt(stmt)
		if err != nil {
			return nil, err
		}
	case ast.WhileStmt:
		return nil, NotImplementedErr
	case ast.BlockStmt:
		newEnv := enviroment.NewEnv(&currentEnv)
		_, err := computeBLockStmt(stmt, &newEnv)
		if err != nil {
			return nil, err
		}
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
	return nil, nil
}

func Interpret(stmts []ast.Stmt) (any, error) {
	for _, stmt := range stmts {
		val, err := execute(stmt)
		if err != nil {
			return nil, err
		}
		if val != nil {
			//return val, nil
			fmt.Println(val)
		}
	}
	return nil, nil
}

func computeIfStmt(stmt ast.Stmt) error {
	condition, err := computeExpr(stmt.Condition)
	if err != nil {
		return err
	}
	if isTruthy(condition) {
		_, err := execute(*stmt.ThenBranch)
		if err != nil {
			return err
		}
	} else {
		_, err := execute(*stmt.ElseBranch)
		if err != nil {
			return err
		}
	}
	return nil
}

// is this the right place ???
// el scope global no tiene parent, de ahi el nil
var currentEnv enviroment.Env = enviroment.NewEnv(nil)

//var prevEnv enviroment.Env = enviroment.NewEnv(nil)

func computeBLockStmt(stmt ast.Stmt, env *enviroment.Env) (any, error) {
	prevNev := currentEnv
	//hago que el env actual sea el que pasa por el argumento
	currentEnv = *env
	val, err := Interpret(stmt.Statements)
	// if err != nil {
	// 	return nil, err
	// }
	//despues de ejecutar el bloque vuelvo a asignar el env a su padre
	currentEnv = prevNev
	return val, err
}

func declareVar(stmt ast.Stmt) error {
	val, err := computeExpr(stmt.Expr)
	if err != nil {
		return err
	}
	err = currentEnv.Define(stmt.Name.Lexeme, val)
	if err != nil {
		return err
	}
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
	case ast.VarAssingmentExpr:
		return computeVarAssingExpr(expr)
	//Estas 4 caen dentro de la cantegoria de operaciones binarias
	case ast.LogicalExpr:
		return nil, NotImplementedErr
	case ast.EqualityExpr, ast.ComparisonExpr, ast.TermExpr, ast.FactorExpr:
		return computeBynaryExpr(expr)
	case ast.UnaryExpr:
		return computeUnaryExpr(expr)
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

func computeVarAssingExpr(expr *ast.Expr) (any, error) {
	val, err := computeExpr(expr.Right) //primero hay que calcular el valor de la expresion que esta a la derecha
	if err != nil {
		return nil, err
	}
	//aqui estoy haciendo que el hecho de asignar una expresion retorne el valor asignado
	//el que esto se haga o no va a depender del paradigma dominante del lenguage
	//por ejemplo los lenguajes funcionales no tienen funciones que no retornen nada
	//en Scala las funciones "vacias" retornan Unit
	return val, currentEnv.Assing(*expr.Name, val)
}
func computeVarExpr(expr *ast.Expr) (any, error) {
	return currentEnv.Get(*expr.Name)

}

func computeUnaryExpr(expr *ast.Expr) (any, error) {
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

func computeBynaryExpr(expr *ast.Expr) (any, error) {

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

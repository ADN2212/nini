package interpreter

import (
	"errors"
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

func execute(stmt ast.Stmt, env *enviroment.Env) (any, error) {
	//fmt.Printf("En env '%v' con valores %v \n", currentEnv.Name, currentEnv.Values)
	switch stmt.Type {
	case ast.VarDeclarationStmt:
		err := declareVar(stmt, env)
		if err != nil {
			return nil, err
		}
	case ast.PrintStmt:
		err := printExpr(stmt.Expr, env)
		if err != nil {
			return nil, err
		}
	case ast.IfStmt:
		err := computeIfStmt(stmt, env)
		if err != nil {
			return nil, err
		}
	case ast.WhileStmt:
		err := computeWhileStmt(stmt, env)
		if err != nil {
			return nil, err
		}
	case ast.BlockStmt:
		_, err := computeBLockStmt(stmt, env)
		if err != nil {
			return nil, err
		}
	case ast.ExprStmt:
		val, err := computeExpr(stmt.Expr, env)
		if err != nil {
			return nil, err
		}
		v, ok := goToNiniRprMap[val]
		if ok {
			val = v
		}
	default:
		return nil, UnknownStmtTypeErr
	}
	return nil, nil
}

func Interpret(stmts []ast.Stmt, env *enviroment.Env) (any, error) {
	if env == nil {
		return nil, errors.New("Env should be defined")
	}
	for _, stmt := range stmts {
		val, err := execute(stmt, env)
		if err != nil {
			return nil, err
		}
		if val != nil {
			fmt.Println(val)
		}
	}
	return nil, nil
}

// tal vez compute no sea el prefijo mas adecuando para esta funciones sobre statements
func computeWhileStmt(stmt ast.Stmt, env *enviroment.Env) error {
	for {
		val, err := computeExpr(stmt.Condition, env)
		if err != nil {
			return err
		}
		if isTruthy(val) {
			//el body de un while se puede interpretar como la then branch de un if isn't it ?
			_, err := execute(*stmt.ThenBranch, env)
			if err != nil {
				return err
			}
		} else {
			break
		}
	}
	return nil
}

func computeIfStmt(stmt ast.Stmt, env *enviroment.Env) error {
	condition, err := computeExpr(stmt.Condition, env)
	if err != nil {
		return err
	}
	if isTruthy(condition) {
		_, err = execute(*stmt.ThenBranch, env)
	}
	if stmt.ElseBranch != nil {
		_, err = execute(*stmt.ElseBranch, env)
	}
	return err
}

func computeBLockStmt(stmt ast.Stmt, env *enviroment.Env) (any, error) {
	newEnv := enviroment.NewEnv("sub env", env)
	val, err := Interpret(stmt.Statements, &newEnv)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func declareVar(stmt ast.Stmt, env *enviroment.Env) error {
	val, err := computeExpr(stmt.Expr, env)
	if err != nil {
		return err
	}
	err = env.Define(stmt.Name.Lexeme, val)
	if err != nil {
		return err
	}
	//fmt.Printf("variable '%v' declarada con valor '%v' \n", stmt.Name.Lexeme, val)
	return nil
}

func printExpr(expr *ast.Expr, env *enviroment.Env) error {
	val, err := computeExpr(expr, env)
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

func computeExpr(expr *ast.Expr, env *enviroment.Env) (any, error) {
	if expr != nil {
		switch expr.Type {
		case ast.VarAssingmentExpr:
			//fmt.Println("In var assigment")
			return computeVarAssingExpr(expr, env)
		//Estas 4 caen dentro de la cantegoria de operaciones binarias
		case ast.LogicalExpr:
			return computeLogicalExpr(*expr, env)
		case ast.EqualityExpr, ast.ComparisonExpr, ast.TermExpr, ast.FactorExpr, ast.ExponetialExpr:
			return computeBynaryExpr(expr, env)
		case ast.UnaryExpr:
			return computeUnaryExpr(expr, env)
		case ast.PrimaryExpr:
			return expr.Val, nil
		case ast.GroupExpr:
			return computeExpr(expr.Right, env)
		case ast.VarExpr:
			return computeVarExpr(expr, env)
		default:
			return nil, UnknowExprTypeErr
		}
	}
	return nil, NilExprErr
}

func computeLogicalExpr(expr ast.Expr, env *enviroment.Env) (any, error) {
	l, err := computeExpr(expr.Left, env)
	if err != nil {
		return nil, err
	}
	leftTrutyVal := isTruthy(l)
	switch expr.Operator.Type {
	//Aqui estoy implementatand short circuiting:
	//Para el operador o; la expresion de la ezquierda es true retornamos true si evaluar la expr de la derecha
	case token.Or:
		if leftTrutyVal {
			//fmt.Println("short cicuiting OR expression")
			return true, nil
		}
		r, err := computeExpr(expr.Right, env)
		if err != nil {
			return nil, err
		}
		return isTruthy(r), nil
	//para el operador y; si la expresion de la izquerda es false, retornamos false sin evaluar la expression de la derecha.
	case token.And:
		if !leftTrutyVal {
			//fmt.Println("short cicuiting AND expression")
			return false, nil
		}
		r, err := computeExpr(expr.Right, env)
		if err != nil {
			return nil, err
		}
		return isTruthy(r), nil
	default:
		return nil, InvalidadLogicalOperator
	}
}

func computeVarAssingExpr(expr *ast.Expr, env *enviroment.Env) (any, error) {
	val, err := computeExpr(expr.Right, env) //primero hay que calcular el valor de la expresion que esta a la derecha
	if err != nil {
		return nil, err
	}
	//aqui estoy haciendo que el hecho de asignar una expresion retorne el valor asignado
	//el que esto se haga o no va a depender del paradigma dominante del lenguage
	//por ejemplo los lenguajes funcionales no tienen funciones que no retornen nada
	//en Scala las funciones "vacias" retornan Unit
	return val, env.Assing(*expr.Name, val)
}
func computeVarExpr(expr *ast.Expr, env *enviroment.Env) (any, error) {
	return env.Get(*expr.Name)
}

func computeUnaryExpr(expr *ast.Expr, env *enviroment.Env) (any, error) {
	r, err := computeExpr(expr.Right, env)
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
	case token.Not:
		return !isTruthy(r), nil
	default:
		return nil, UnknowUnaryOperator
	}
}

func computeBynaryExpr(expr *ast.Expr, env *enviroment.Env) (any, error) {

	l, err := computeExpr(expr.Left, env)
	if err != nil {
		return 0, err
	}
	r, err := computeExpr(expr.Right, env)
	if err != nil {
		return 0, err
	}

	leftVal, okL := l.(float64)
	rightVal, okR := r.(float64)

	//Ojo: todo lo que esta dentro de los if pueden ser funciones helpers para minimizar la logitud de computeBynaryExpr.
	if okL && okR {
		return operateNumbers(leftVal, rightVal, expr.Operator)
	}

	leftStr, okL := l.(string)
	rightStr, okR := r.(string)
	if okL && okR {
		return operateStrings(leftStr, rightStr, expr.Operator)
	}

	leftBool, okL := l.(bool)
	rightBool, okR := r.(bool)
	if okL && okR {
		return operateBooleans(leftBool, rightBool, expr.Operator)
	}

	leftStr, okL = l.(string)
	rightNum, okR := r.(float64)
	if okL && okR && expr.Operator.Type == token.Plus {
		return concatenateStrAndNum(leftStr, rightNum, false), nil
	}

	leftNum, okL := l.(float64)
	rightStr, okR = r.(string)
	if okL && okR && expr.Operator.Type == token.Plus {
		return concatenateStrAndNum(rightStr, leftNum, true), nil
	}

	return nil, UnknowBynaryOperation
}

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
	//el ordern importa aqui ???
	case ast.VarDeclarationStmt:
		//fmt.Println("Por aqui?")
		err := declareVar(stmt, env)
		if err != nil {
			return nil, err
		}
	case ast.PrintStmt:
		//fmt.Println("Por aca ??")
		err := printExpr(stmt.Expr, env)
		if err != nil {
			return nil, err
		}
	case ast.IfStmt:
		//fmt.Println("Por acuya ???")
		err := computeIfStmt(stmt, env)
		if err != nil {
			return nil, err
		}
	case ast.WhileStmt:
		//fmt.Println("In while stmt ...")
		err := computeWhileStmt(stmt, env)
		if err != nil {
			return nil, err
		}
		//return nil, NotImplementedErr
	case ast.BlockStmt:
		//fmt.Println("In block statement ...")
		//newEnv := enviroment.NewEnv(&currentEnv)
		_, err := computeBLockStmt(stmt, env)
		if err != nil {
			return nil, err
		}
		//fmt.Println(env.Name)
	case ast.ExprStmt:
		//fmt.Println("Expr statement")
		val, err := computeExpr(stmt.Expr, env)
		if err != nil {
			return nil, err
		}
		v, ok := goToNiniRprMap[val]
		if ok {
			val = v
		}
		//fmt.Printf("expresion ejecutada => %v \n", val)
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
		//fmt.Println("Starting iteration ...")
		val, err := computeExpr(stmt.Condition, env)
		if err != nil {
			//fmt.Println("Aqui ??")
			return err
		}
		if isTruthy(val) {
			//el body de un while se puede interpretar como la then branch de un if
			//isn't it ?
			_, err := execute(*stmt.ThenBranch, env)
			if err != nil {
				//fmt.Println("o aqui ??")
				return err
			}
		} else {
			break
		}
	}
	return nil
}

func computeIfStmt(stmt ast.Stmt, env *enviroment.Env) error {
	//fmt.Println("In if statement")
	condition, err := computeExpr(stmt.Condition, env)
	if err != nil {
		return err
	}
	if isTruthy(condition) {
		_, err := execute(*stmt.ThenBranch, env)
		if err != nil {
			return err
		}
	} else {
		_, err := execute(*stmt.ElseBranch, env)
		if err != nil {
			return err
		}
	}
	return nil
}

// is this the right place ???
// el scope global no tiene parent, de ahi el nil
//var currentEnv enviroment.Env = enviroment.NewEnv("global", nil)

//var prevEnv enviroment.Env = enviroment.NewEnv(nil)

func computeBLockStmt(stmt ast.Stmt, env *enviroment.Env) (any, error) {
	//fmt.Println("Entrando al bloque")
	newEnv := enviroment.NewEnv("sub env", env)
	//fmt.Println(newEnv.ParentEnv.Values)
	val, err := Interpret(stmt.Statements, &newEnv)
	if err != nil {
		//fmt.Println("Here")
		return nil, err
	}
	//fmt.Println("Saliendo del bloque")
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
	//fmt.Println("Computing expr ...")
	if expr != nil {
		switch expr.Type {
		case ast.VarAssingmentExpr:
			//fmt.Println("In var assigment")
			return computeVarAssingExpr(expr, env)
		//Estas 4 caen dentro de la cantegoria de operaciones binarias
		case ast.LogicalExpr:
			//fmt.Println(2)
			//return nil, NotImplementedErr
			return computeLogicalExpr(*expr, env)
		case ast.EqualityExpr, ast.ComparisonExpr, ast.TermExpr, ast.FactorExpr:
			//fmt.Println("Tiene que ser por aqui ???")
			return computeBynaryExpr(expr, env)
		case ast.UnaryExpr:
			//fmt.Println(3)
			return computeUnaryExpr(expr, env)
		case ast.PrimaryExpr:
			//fmt.Println(4)
			return expr.Val, nil
		case ast.GroupExpr:
			//fmt.Println(5)
			return computeExpr(expr.Right, env)
		case ast.VarExpr:
			//fmt.Println(6)
			return computeVarExpr(expr, env)
		default:
			return nil, UnknowExprTypeErr
		}
	}
	return nil, NilExprErr
}

func computeLogicalExpr(expr ast.Expr, env *enviroment.Env) (any, error) {
	//fmt.Println("In logical exprds")
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
			fmt.Println("short cicuiting OR expression")
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
			fmt.Println("short cicuiting AND expression")
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
	//fmt.Println(currentEnv.Name)
	//fmt.Println(currentEnv.ParentEnv == &currentEnv)
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
	case token.Not: //esto es el no
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

package interpreter

import (
	"errors"
	"fmt"
	"github.com/ADN2212/nini/ast"
	"github.com/ADN2212/nini/enviroment"
	"github.com/ADN2212/nini/token"
)

type interpreter struct {
	stmts   []ast.Stmt
	err     error
	env     *enviroment.Env
	current int
}

func NewInterpreter(stmts []ast.Stmt, env *enviroment.Env) interpreter {
	return interpreter{
		stmts: stmts,
		env:   env,
	}
}

func (i *interpreter) getCurrentStmt() ast.Stmt {
	return i.stmts[i.current] //esto puede paniquear
}

func (i *interpreter) hasNext() bool {
	return i.current < len(i.stmts)
}

func (i *interpreter) Interpret() (any, error) {

	if i.env == nil {
		return nil, errors.New("Env should be defined")
	}

	for i.hasNext() {
		i.execute(nil)
		if i.err != nil {
			return nil, i.err
		}
	}

	return nil, nil
}

func (i *interpreter) execute(s *ast.Stmt) {
	//fmt.Printf("En env '%v' con valores %v \n", currentEnv.Name, currentEnv.Values)
	//esto es solo un hack asqueroso, plis fix this.
	var stmt ast.Stmt
	if s == nil {
		stmt = i.getCurrentStmt()
	} else {
		stmt = *s
	}
	
	switch stmt.Type {
	case ast.VarDeclarationStmt:
		i.declareVar(&stmt)
	case ast.PrintStmt:
		i.printExpr(stmt.Expr)
	case ast.IfStmt:
		i.executeIfStmt(&stmt)
	case ast.WhileStmt:
		
	case ast.BlockStmt:
		i.executeBLockStmt(&stmt)
	case ast.ExprStmt:
		i.computeExpr(stmt.Expr)
	default:
		i.err = UnknownStmtTypeErr
	}
	i.current++ //iterar al siguiente statement
}

func (i *interpreter) printExpr(expr *ast.Expr) {
	val := i.computeExpr(expr)
	if i.err == nil {
		v, ok := goToNiniRprMap[val]
		if ok {
			val = v
		}
		fmt.Printf("mira => %v \n", val)
	}
}

func (i *interpreter) declareVar(stmt *ast.Stmt) {
	val := i.computeExpr(stmt.Expr)
	err := i.env.Define(stmt.Name.Lexeme, val)
	if err != nil {
		i.err = err
	}
}

func (i *interpreter) computeExpr(expr *ast.Expr) any {
	if expr != nil {
		switch expr.Type {
		case ast.PrimaryExpr:
			return expr.Val
		case ast.VarAssingmentExpr:
			return i.computeVarAssingExpr(expr)
		//Estas 4 caen dentro de la cantegoria de operaciones binarias
		case ast.LogicalExpr:
			return i.computeLogicalExpr(expr)
		case ast.EqualityExpr, ast.ComparisonExpr, ast.TermExpr, ast.FactorExpr, ast.ExponetialExpr:
			return i.computeBynaryExpr(expr)
		case ast.UnaryExpr:
			return i.computeUnaryExpr(expr)
		case ast.GroupExpr:
			return i.computeExpr(expr.Right)
		case ast.VarExpr:
			return i.computeVarExpr(expr)
		default:
			i.err = UnknowExprTypeErr
			return nil
		}
	}
	i.err = NilExprErr
	return nil
}

func (i *interpreter) computeVarAssingExpr(expr *ast.Expr) any {
	val := i.computeExpr(expr.Right) //primero hay que calcular el valor de la expresion que esta a la derecha
	err := i.env.Assing(*expr.Name, val)

	if err != nil {
		i.err = err
		return nil
	}
	//aqui estoy haciendo que el hecho de asignar una expresion retorne el valor asignado
	//el que esto se haga o no va a depender del paradigma dominante del lenguage
	//por ejemplo los lenguajes funcionales no tienen funciones que no retornen nada
	//en Scala las funciones "vacias" retornan Unit
	//return val, i.env.Assing(*expr.Name, val)
	return val
}

// tal vez compute no sea el prefijo mas adecuando para esta funciones sobre statements
// func computeWhileStmt(stmt ast.Stmt, env *enviroment.Env) error {
// 	for {
// 		val, err := computeExpr(stmt.Condition, env)
// 		if err != nil {
// 			return err
// 		}
// 		if isTruthy(val) {
// 			//el body de un while se puede interpretar como la then branch de un if isn't it ?
// 			_, err := execute(*stmt.ThenBranch, env)
// 			if err != nil {
// 				return err
// 			}
// 		} else {
// 			break
// 		}
// 	}
// 	return nil
// }

func (i *interpreter) executeIfStmt(stmt *ast.Stmt) {
	condition := i.computeExpr(stmt.Condition)
	if i.err != nil {
		return //esto es solo para detener la ejecucion en caso de que pase un error
	}
	if isTruthy(condition) {
		i.execute(stmt.ThenBranch)
		return//sin este return se ejecutaria la esle branch en caso de existir.
	}
	if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}
}

func (i *interpreter) executeBLockStmt(stmt *ast.Stmt) (any, error) {
	newEnv := enviroment.NewEnv("sub env", i.env)
	//como newEnv tiene como padre el env del interprete que llama a este metodo,
	//aunque sea otra instancai podra acceder a las var del interprete anterior.
	newInterpreter := NewInterpreter(stmt.Statements, &newEnv)
	val, err := newInterpreter.Interpret()
	if err != nil {
		return nil, err
	}
	return val, nil//esto tiene que retornar algo ???
}

func (i *interpreter) computeLogicalExpr(expr *ast.Expr) any {
	l := i.computeExpr(expr.Left)

	if i.err != nil {
		return nil
	}

	leftTrutyVal := isTruthy(l)
	switch expr.Operator.Type {
	//Aqui estoy implementatand short circuiting:
	//Para el operador o; la expresion de la ezquierda es true retornamos true si evaluar la expr de la derecha
	case token.Or:
		if leftTrutyVal {
			return true
		}
		r := i.computeExpr(expr.Right)
		if i.err != nil {
			return nil
		}
		return isTruthy(r)
	//para el operador y; si la expresion de la izquerda es false, retornamos false sin evaluar la expression de la derecha.
	case token.And:
		if !leftTrutyVal {
			return false
		}
		r := i.computeExpr(expr.Right)
		if i.err != nil {
			return nil
		}
		return isTruthy(r)
	default:
		i.err = InvalidadLogicalOperator
		return nil
	}
}

func (i *interpreter) computeVarExpr(expr *ast.Expr) any {
	val, err := i.env.Get(*expr.Name)
	if err != nil {
		i.err = err
		return nil
	}
	return val
}

func (i *interpreter) computeUnaryExpr(expr *ast.Expr) any {
	r := i.computeExpr(expr.Right)
	if i.err != nil {
		return nil
	}
	switch expr.Operator.Type {
	case token.Minus:
		num, ok := r.(float64)
		if !ok {
			i.err = ShouldBeNumberErr
			return nil
		}
		return -num
	case token.Not:
		return !isTruthy(r)
	default:
		i.err = UnknowUnaryOperator
		return nil
	}
}

func (i *interpreter) computeBynaryExpr(expr *ast.Expr) any {

	l := i.computeExpr(expr.Left)
	//igual estos early returns seran utiles para evitar llegar operaciones con esados invalidos
	//does it make sence??
	if i.err != nil {
		return nil
	}
	r := i.computeExpr(expr.Right)
	if i.err != nil {
		return nil
	}

	//fmt.Printf("operando %v %v %v \n", l, expr.Operator.Lexeme, r)

	leftVal, okL := l.(float64)
	rightVal, okR := r.(float64)

	//Ojo: todo lo que esta dentro de los if pueden ser funciones helpers para minimizar la logitud de computeBynaryExpr.
	if okL && okR {
		val, err := operateNumbers(leftVal, rightVal, expr.Operator)
		//este patron se va a repetir
		if err != nil {
			i.err = err
			return nil
		}
		return val
	}

	leftStr, okL := l.(string)
	rightStr, okR := r.(string)
	if okL && okR {
		val, err := operateStrings(leftStr, rightStr, expr.Operator)
		if err != nil {
			i.err = err
			return nil
		}
		return val
	}

	leftBool, okL := l.(bool)
	rightBool, okR := r.(bool)
	if okL && okR {
		val, err := operateBooleans(leftBool, rightBool, expr.Operator)
		if err != nil {
			i.err = err
			return nil
		}
		return val
	}

	leftStr, okL = l.(string)
	rightNum, okR := r.(float64)
	if okL && okR && expr.Operator.Type == token.Plus {
		return concatenateStrAndNum(leftStr, rightNum, false)
	}

	leftNum, okL := l.(float64)
	rightStr, okR = r.(string)
	if okL && okR && expr.Operator.Type == token.Plus {
		return concatenateStrAndNum(rightStr, leftNum, true)
	}

	i.err = UnknowBynaryOperation
	return nil
}

// esto es para fines de representacion, osea, lo que el user del lenguage ve.
var goToNiniRprMap map[any]string = map[any]string{
	false: "falso",
	true:  "verdadero",
	nil:   "nulo",
}

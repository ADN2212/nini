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

// func (i *interpreter) getCurrentStmtType() ast.StmtType {
// 	return i.stmts[i.current].Type
// }

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
		i.execute()
		if i.err != nil {
			return nil, i.err
		}
	}
	return nil, nil
}

func (i *interpreter) execute() {
	//fmt.Printf("En env '%v' con valores %v \n", currentEnv.Name, currentEnv.Values)
	stmt := i.getCurrentStmt()
	switch stmt.Type {
	case ast.VarDeclarationStmt:
		fmt.Println("In var declaration statement")
		i.declareVar(&stmt)
	case ast.PrintStmt:
		fmt.Println("In print statement")
		i.printExpr(stmt.Expr)
		// err := printExpr(stmt.Expr, env)
		// if err != nil {
		// 	return nil, err
		// }
	case ast.IfStmt:
		// err := computeIfStmt(stmt, env)
		// if err != nil {
		// 	return nil, err
		// }
	case ast.WhileStmt:
		// err := computeWhileStmt(stmt, env)
		// if err != nil {
		// 	return nil, err
		// }
	case ast.BlockStmt:
		// _, err := computeBLockStmt(stmt, env)
		// if err != nil {
		// 	return nil, err
		// }
	case ast.ExprStmt:
		fmt.Println("In expression statement")
		val := i.computeExpr(stmt.Expr)
		v, ok := goToNiniRprMap[val]
		if ok {
			val = v
		}
		fmt.Printf("expression returned: %v \n", val)
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
		//fmt.Println("Aqui ??")
		i.err = err
	}
	//fmt.Println(i.env)
	fmt.Printf("variable '%v' declarada con valor '%v' \n", stmt.Name.Lexeme, val)
	//return nil
}

func (i *interpreter) computeExpr(expr *ast.Expr) any {

	//esto puede ser otro metodo
	//exprType := i.getCurrentStmt().Expr.Type

	if expr != nil {
		switch expr.Type {
		case ast.PrimaryExpr:
			//fmt.Println("In primary expression")
			return expr.Val
		case ast.VarAssingmentExpr:
			//fmt.Println("In var assigment expresssion")
			return i.computeVarAssingExpr(expr)
		//Estas 4 caen dentro de la cantegoria de operaciones binarias
		case ast.LogicalExpr:
			fmt.Println("Computing logical expression")
			return i.computeLogicalExpr(expr)
		case ast.EqualityExpr, ast.ComparisonExpr, ast.TermExpr, ast.FactorExpr, ast.ExponetialExpr:
			fmt.Println("In bynary expression")
			return i.computeBynaryExpr(expr)
		case ast.UnaryExpr:
			fmt.Println("In unary expr")
			return i.computeUnaryExpr(expr)
		// case ast.GroupExpr:
		// 	return computeExpr(expr.Right, env)
		case ast.VarExpr:
			return i.computeVarExpr(expr)
		default:
			//return nil, UnknowExprTypeErr
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

// func computeIfStmt(stmt ast.Stmt, env *enviroment.Env) error {
// 	condition, err := computeExpr(stmt.Condition, env)
// 	if err != nil {
// 		return err
// 	}
// 	if isTruthy(condition) {
// 		_, err = execute(*stmt.ThenBranch, env)
// 	}
// 	if stmt.ElseBranch != nil {
// 		_, err = execute(*stmt.ElseBranch, env)
// 	}
// 	return err
// }

// func computeBLockStmt(stmt ast.Stmt, env *enviroment.Env) (any, error) {
// 	newEnv := enviroment.NewEnv("sub env", env)
// 	val, err := Interpret(stmt.Statements, &newEnv)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return val, nil
// }

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
			//fmt.Println("short cicuiting OR expression")
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
			//fmt.Println("short cicuiting AND expression")
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

	fmt.Printf("operando %v %v %v \n", l, expr.Operator.Lexeme, r)

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

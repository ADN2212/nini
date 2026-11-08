package main

import (
	"bufio"
	"fmt"
	"github.com/ADN2212/nini/lexer"
	"github.com/ADN2212/nini/parser"
	"os"
)

func main() {
	argsNum := len(os.Args)
	switch argsNum {
	case 2:
		fmt.Println("Ejecutando script")
	case 1:
		fmt.Println("Starting commant promtp")
		runPrompt()
	default:
		fmt.Println("Opcion no definida")
	}
}

// Esta es la funcion que se usa cuando el user esta usando el REPL
func runPrompt() {
	//fmt.Println("running prompt ...")
	input := bufio.NewScanner(os.Stdin) //Este es el scanner para lo que escirba el user
	for {
		fmt.Print("#nini>>")
		ok := input.Scan() //escanea la linea
		if !ok {
			break
		}
		line := input.Text() //estrae la linea
		if input.Err() != nil {
			os.Exit(64)
		}
		//fmt.Println(line)
		err := run(line)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(64)
		}
	}
}

func run(src string) error {
	tokens, err := lexer.ScanTokens(src)
	if err != nil {
		return err
	}
	for _, t := range tokens {
		fmt.Printf("lexema: %v, tipo: %v, valor: %v\n", t.Lexeme, t.Type, t.Literal)
	}

	stmts, err := parser.Parse(tokens)
	if err != nil {
		return err
	}

	expr := stmts[0].Expr
	
	//fmt.Println(expr)
	//fmt.Println(expr.Operator)
	//fmt.Println(expr.Left.Val)
	
	//fmt.Println(expr.Left)

	//calls := 0

	for {
		fmt.Println(expr.Operator)
		//fmt.Println(expr == expr.Right)
		//fmt.Println(expr == expr.Left)
		expr = expr.Right
		if expr == nil {
			//fmt.Println(expr.Operator)
			fmt.Println("Saliendo")
			return nil
		}
		fmt.Println(expr.Right)
		//calls++
	}

	//fmt.Println(expr.Operator)

	// for _, stmt := range stmts {
	// 	fmt.Println(stmt.Type)
	// 	fmt.Println(stmt.Expr.Operator)
	// 	fmt.Println(stmt.Expr.Left.Left.Left.Left.Left)
	// 	//fmt.Println(stmt.Expr.Left.Left.Operator)
	// 	//fmt.Println(stmt.Expr.Left.Right.Val)

	// }

	return nil
}

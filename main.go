package main

import (
	"bufio"
	"fmt"
	"github.com/ADN2212/nini/lexer"
	"os"
	"github.com/ADN2212/nini/parser"
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
	

	for _, stmt := range stmts {
		fmt.Println(stmt)
	}

	return nil	
}

package main

import (
	"bufio"
	"fmt"
	"github.com/ADN2212/nini/enviroment"
	"github.com/ADN2212/nini/interpreter"
	"github.com/ADN2212/nini/lexer"
	"github.com/ADN2212/nini/parser"
	"os"
)

func main() {
	argsNum := len(os.Args)
	globalEnv := enviroment.NewEnv("global", nil)
	switch argsNum {
	case 2:
		fmt.Println("Ejecutando script")
		runFile(os.Args[1], &globalEnv)
	case 1:
		//globalEnv := enviroment.NewEnv("global", nil)
		fmt.Println("Starting commant promtp")
		runPrompt(&globalEnv)
	default:
		fmt.Println("Opcion no definida")
	}
}

func runFile(path string, env *enviroment.Env) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(64)
	}
	err = run(string(bytes), env)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(64)
	}
}

// Esta es la funcion que se usa cuando el user esta usando el REPL
func runPrompt(env *enviroment.Env) {
	//fmt.Println("running prompt ...")
	input := bufio.NewScanner(os.Stdin) //Este es el scanner para lo que escirba el user
	//globalEnv := enviroment.NewEnv("global", nil)
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
		err := run(line, env)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(64)
		}
	}
}

func run(src string, env *enviroment.Env) error {
	tokens, err := lexer.ScanTokens(src)
	if err != nil {
		return err
	}

	stmts, err := parser.Parse(tokens)
	if err != nil {
		return err
	}

	res, err := interpreter.Interpret(stmts, env)
	if err != nil {
		return err
	}
	if res != nil {
		fmt.Println(res)
	}
	return nil
}

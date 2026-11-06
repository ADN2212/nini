package main

import (
	"fmt"
	"os"
)

func main() {
	argsNum := len(os.Args)
	switch argsNum {
	case 2:
		fmt.Println("Ejecutando script")
	case 1:
		fmt.Println("Starting commant promtp")
	default:
		fmt.Println("Opcion no definida")
	}
}

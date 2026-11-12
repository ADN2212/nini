package enviroment

import (
	"fmt"

	"github.com/ADN2212/nini/token"
)

// hay alguna ventaja en declararlo asi a que sea un simple map nombrado como tipo???
type Env struct {
	values map[string]any
}

func NewEnv() Env {
	return Env{values: map[string]any{}}
}

func (env *Env) Define(name string, value any) error {
	_, ok := env.values[name]
	if ok {
		return AlreadyDefinedVar
	}
	env.values[name] = value
	return nil
}

func (env *Env) Get(name token.Token) (any, error) {
	val, ok := env.values[name.Lexeme]
	if !ok {
		return nil, UndefinedVarError
	}
	return val, nil
}

func (env *Env) Assing(name token.Token, value any) error {
	_, ok := env.values[name.Lexeme]
	//si no se hiciera esto el nini tendria declaracion de variables implicita
	if !ok {
		return UndefinedVarError
	}
	env.values[name.Lexeme] = value
	fmt.Printf("la variable '%v' ha sido reasignada al valor '%v' \n", name.Lexeme, value)
	return nil
}

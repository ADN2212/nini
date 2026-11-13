package enviroment

import (
	"fmt"
	"github.com/ADN2212/nini/token"
)

// hay alguna ventaja en declararlo asi a que sea un simple map nombrado como tipo???
type Env struct {
	//en el libor a esto le llaman enclosing
	parentEnv *Env
	values    map[string]any
}

func NewEnv(parent *Env) Env {
	return Env{parentEnv: parent, values: map[string]any{}}
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
		//si la variable no aparece en el scope actual
		//la buscamos en el scope padre en caso de que este exista
		if env.parentEnv != nil {
			return env.parentEnv.Get(name)
		}
		//de lo contrario podemos concluir que la perra no existe.
		return nil, UndefinedVarError
	}
	return val, nil
}

func (env *Env) Assing(name token.Token, value any) error {
	_, ok := env.values[name.Lexeme]
	//si no se hiciera esto nini tendria declaracion de variables implicita
	if !ok {
		//el mismo principio que en el Get, si la variable no existen en este scope busca en el padre recursivamente
		if env.parentEnv != nil {
			env.parentEnv.Assing(name, value)
		}
		return UndefinedVarError
	}
	env.values[name.Lexeme] = value
	fmt.Printf("la variable '%v' ha sido reasignada al valor '%v' \n", name.Lexeme, value)
	return nil
}

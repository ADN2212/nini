package enviroment

import (
	"fmt"
	"github.com/ADN2212/nini/token"
)

// hay alguna ventaja en declararlo asi a que sea un simple map nombrado como tipo???
type Env struct {
	//en el libor a esto le llaman enclosing
	Name      string
	ParentEnv *Env
	Values    map[string]any
}

func NewEnv(name string, parent *Env) Env {
	return Env{Name: name, ParentEnv: parent, Values: map[string]any{}}
}

func (env *Env) Define(name string, value any) error {
	_, ok := env.Values[name]
	if ok {
		return AlreadyDefinedVar
	}
	env.Values[name] = value
	//fmt.Println(env.Values)
	return nil
}

func (env *Env) Get(name token.Token) (any, error) {
	//fmt.Printf("Calling Get for %v \n", env.Name)
	//fmt.Println(env.Values)
	val, ok := env.Values[name.Lexeme]
	if !ok {
		//si la variable no aparece en el scope actual
		//la buscamos en el scope padre en caso de que este exista
		if env.ParentEnv != nil {
			return env.ParentEnv.Get(name)
		}
		//de lo contrario podemos concluir que la perra no existe.
		//fmt.Println("En Get")
		return nil, UndefinedVarError
	}
	return val, nil
}

func (env *Env) Assing(name token.Token, value any) error {
	_, ok := env.Values[name.Lexeme]
	//si no se hiciera esto nini tendria declaracion de variables implicita
	if !ok {
		//el mismo principio que en el Get, si la variable no existen en este scope busca en el padre recursivamente
		if env.ParentEnv != nil {
			env.ParentEnv.Assing(name, value)
			return nil
		}
		fmt.Printf("En Assing para %v \n", env.Name)
		return UndefinedVarError
	}
	env.Values[name.Lexeme] = value
	//fmt.Println(env.Values)
	//fmt.Printf("la variable '%v' ha sido reasignada al valor '%v' \n", name.Lexeme, value)
	return nil
}

// func (env *Env) Copy() Env {
// 	return Env{
// 		Name:      env.Name,
// 		ParentEnv: env.ParentEnv,
// 		Values:    env.Values,
// 	}
// }

// func (env *Env) Set(e *Env) {
// 	env.Name = e.Name
// 	env.ParentEnv = e.ParentEnv
// 	env.Values = e.Values
// }

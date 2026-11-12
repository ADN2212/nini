package enviroment

import (
	"errors"

	"github.com/ADN2212/nini/token"
)

// hay alguna ventaja en declararlo asi a que sea un simple map nombrado como tipo???
type Env struct {
	values map[string]any
}

func NewEnv() Env {
	return Env{values: map[string]any{}}
}

// por la forma en la que esto esta escrito
// se puede hacer lo siguiente
// var x = 10;
// var x = 25;
// es decir, redeclarar variables
// eso puede ser una decicion de disenio, o se puede agragar un constrain
func (env *Env) Define(name string, value any) {
	env.values[name] = value
}

// esto es un runtime error
var UndefinedVarError error = errors.New("Variable no definida")

func (env *Env) Get(name token.Token) (any, error) {
	val, ok := env.values[name.Lexeme]
	if !ok {
		return nil, UndefinedVarError
	}
	return val, nil
}

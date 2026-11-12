package enviroment

import "errors"

// Estos errores sucenden en el interpreter, so they are runtime errors
var AlreadyDefinedVar error = errors.New("Esta variable ya ha sido declarada")
var UndefinedVarError error = errors.New("Variable no definida")

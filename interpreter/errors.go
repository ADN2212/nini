package interpreter

import "errors"

// mover esto a un file a parte
// var NotImplementedErr error = errors.New("Esta vaina no ha sido implementada")
var UnknownStmtTypeErr error = errors.New("Tipo de statement desconocido")
var UnknowExprTypeErr error = errors.New("Tipo de expression desconocido")
var ShouldBeNumberErr error = errors.New("Esto deberia ser un numero")
var UnknowUnaryOperator error = errors.New("Operador unario desconocido")
var UnknowBynaryOperator error = errors.New("Operador binario desconocido")
var UnknowBynaryOperation error = errors.New("Operacion binaria desconocida")
var InvalidadLogicalOperator error = errors.New("Operador logico invalido")
var NilExprErr error = errors.New("Expression nula no calculable")
var ModOnlyForIntsErr error = errors.New("el operador '%' solo esta definido para enteros")

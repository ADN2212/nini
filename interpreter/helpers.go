package interpreter

import (
	"fmt"
	"github.com/ADN2212/nini/token"
)

func operateNumbers(a, b float64, operator *token.Token) (any, error) {
	switch operator.Type {
	//aritmetica
	case token.Plus:
		return a + b, nil
	case token.Minus:
		return a - b, nil
	case token.Slash:
		return a / b, nil
	case token.Star:
		return a * b, nil
	case token.Mod:
		aInt := int(a)
		bInt := int(b)
		if a-float64(aInt) != 0 || b-float64(bInt) != 0 {
			return nil, ModOnlyForIntsErr
		} else {
			//como en nini todo numero es un float64,
			// debo parsear el resultado para que se puedan hacer operaciones conjuntas
			return float64(aInt % bInt), nil
		}
	//comparaciones:
	case token.Greater:
		return a > b, nil
	case token.GreaterEqual:
		return a >= b, nil
	case token.Less:
		return a < b, nil
	case token.LessEqual:
		return a <= b, nil

	//Igualdades:
	case token.BangEqual: //!=
		return a != b, nil
	case token.EqualEqual:
		return a == b, nil
	default:
		return nil, UnknowBynaryOperator
	}
}

func operateStrings(a, b string, operator *token.Token) (any, error) {
	switch operator.Type {
	case token.Plus:
		return a + b, nil
	case token.BangEqual:
		return a != b, nil
	case token.EqualEqual:
		return a == b, nil
	default:
		return nil, UnknowBynaryOperator
	}
}

func operateBooleans(a, b bool, operator *token.Token) (bool, error) {
	switch operator.Type {
	case token.BangEqual:
		return a != b, nil
	case token.EqualEqual:
		return a == b, nil
	}
	return false, UnknowBynaryOperator
}

// this is also a helper
func concatenateStrAndNum(str string, num float64, flip bool) string {
	if flip {
		return fmt.Sprintf("%v%v", num, str)
	} else {
		return fmt.Sprintf("%v%v", str, num)
	}
}

// this is a helper
// ojo, la truthyness so podria definir de forma diferente.
func isTruthy(val any) bool {
	if val == nil {
		return false
	}
	b, ok := val.(bool)
	if ok {
		return b
	}
	return true
}

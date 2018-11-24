package evaluator

import (
	"monkey/ast"
	"monkey/object"
	"fmt"
)
var (
	NULL = &object.Null{}
	TRUE = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)
func Eval(node ast.Node) object.Object {
	switch node:= node.(type) {
	case *ast.Program:
		return EvalProgram(node)
	case *ast.BlockStatement:
		return EvalBlockStatements(node)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node)

	case *ast.IntegerLiteral:
		return &object.Integer{ Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBoolewnObject(node.Value)
	}
	return nil
}
func EvalProgram(program *ast.Program) object.Object {
	var result object.Object 

	for _,statement := range program.Statements {
		result = Eval(statement)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}	
	return result 
}

func EvalBlockStatements(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement)

		if result != nil {

			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

func nativeBoolToBoolewnObject(input bool) *object.Boolean { if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinuxPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinuxPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBoolewnObject(left == right)
	case operator == "!=":
		return nativeBoolToBoolewnObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",left.Type(),operator,right.Type())
	default:
		return newError("unknown operator: %s %s %s",left.Type(),operator,right.Type())
	}
}
func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftval := left.(*object.Integer).Value
	rightval := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftval + rightval}
	case "-":
		return &object.Integer{Value: leftval - rightval}
	case "*":
		return &object.Integer{Value: leftval * rightval}
	case "/":
		return &object.Integer{Value: leftval / rightval}
	case "<":
		return nativeBoolToBoolewnObject(leftval < rightval)
	case ">":
		return nativeBoolToBoolewnObject(leftval > rightval)
	case "==":
		return nativeBoolToBoolewnObject(leftval == rightval)
	case "!=":
		return nativeBoolToBoolewnObject(leftval != rightval)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}
func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)

	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj{
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
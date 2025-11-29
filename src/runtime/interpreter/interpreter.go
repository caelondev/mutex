package interpreter

import (
	"fmt"
	"math"

	"github.com/caelondev/mutex/src/frontend/ast"
	"github.com/caelondev/mutex/src/frontend/lexer"
)

// EvaluateStatement evaluates a statement node
func EvaluateStatement(node ast.Statement) RuntimeValue {
	switch n := node.(type) {
	case *ast.BlockStatement:
		return evaluateBlockStatement(n)
	case *ast.ExpressionStatement:
		return EvaluateExpression(n.Expression)
	default:
		panic(fmt.Sprintf("Unsupported statement node type: %T", node))
	}
}

// EvaluateExpression evaluates an expression node
func EvaluateExpression(node ast.Expression) RuntimeValue {
	switch n := node.(type) {
	case *ast.NumberExpression:
		return evaluateNumberExpression(n)
	case *ast.BinaryExpression:
		return evaluateBinaryExpression(n)
	default:
		panic(fmt.Sprintf("Unsupported expression node type: %T", node))
	}
}

func evaluateBlockStatement(block *ast.BlockStatement) RuntimeValue {
	var lastEvaluated RuntimeValue = &NullValue{}
	
	for _, statement := range block.Body {
		lastEvaluated = EvaluateStatement(statement)
	}
	
	return lastEvaluated
}

func evaluateNumberExpression(expr *ast.NumberExpression) RuntimeValue {
	return &NumberValue{Value: expr.Value}
}

func evaluateBinaryExpression(expr *ast.BinaryExpression) RuntimeValue {
	left := EvaluateExpression(expr.Left)
	right := EvaluateExpression(expr.Right)
	
	// Ensure both sides are numbers
	leftNum, leftOk := left.(*NumberValue)
	rightNum, rightOk := right.(*NumberValue)
	
	if !leftOk || !rightOk {
		panic("Binary operations require numeric operands")
	}
	
	return evaluateNumericBinaryExpression(leftNum, rightNum, expr.Operator)
}

func evaluateNumericBinaryExpression(left *NumberValue, right *NumberValue, operator lexer.Token) RuntimeValue {
	result := 0.0
	
	switch operator.TokenType {
	case lexer.PLUS:
		result = left.Value + right.Value
	case lexer.MINUS:
		result = left.Value - right.Value
	case lexer.STAR:
		result = left.Value * right.Value
	case lexer.SLASH:
		if right.Value == 0 {
			panic("Division by zero")
		}
		result = left.Value / right.Value
	case lexer.MODULO:
		if right.Value == 0 {
			panic("Modulo by zero")
		}
		result = math.Mod(left.Value, right.Value)
	default:
		panic(fmt.Sprintf("Unsupported binary operator: %s", operator.Lexeme))
	}
	
	return &NumberValue{Value: result}
}

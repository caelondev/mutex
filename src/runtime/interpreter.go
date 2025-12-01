package runtime

import (
	"fmt"
	"os"

	"github.com/caelondev/mutex/src/frontend/ast"
	"github.com/sanity-io/litter"
)

func EvaluateExpression(node ast.Expression, env Environment) RuntimeValue {
	switch n := node.(type) {
	case *ast.NumberExpression:
		return evaluateNumberExpression(n)
	case *ast.StringExpression:
		return evaluateStringExpression(n)
	case *ast.BinaryExpression:
		return evaluateBinaryExpression(n, env)
	case *ast.SymbolExpression:
		return evaluateSymbolExpression(n, env)
	case *ast.AssignmentExpression:
		return evaluateAssignmentExpression(n, env)
	case *ast.UnaryExpression:
		return evaluateUnaryExpression(n, env)
	case *ast.PostfixExpression:
		return evaluatePostfixExpression(n, env)
	case *ast.ArrayExpression:
		return evaluateArrayExpression(n, env)
	case *ast.ArrayIndexExpression:
		return evaluateIndexExpression(n, env)
	case *ast.ArrayIndexAssignmentExpression:
		return evaluateIndexAssignmentExpression(n, env)
	case *ast.CallExpression:
		return evaluateCallExpression(n, env)

	default:
		litter.Dump(fmt.Sprintf("Unsupported expression node type %v\n", node))
		litter.Dump(node)
		os.Exit(65)
	}

	return nil
}

func EvaluateStatement(node ast.Statement, env Environment) RuntimeValue {
	switch n := node.(type) {
	case *ast.BlockStatement:
		return evaluateBlockStatement(n, env)
	case *ast.ExpressionStatement:
		return EvaluateExpression(n.Expression, env)
	case *ast.VariableDeclarationStatement:
		return evaluateVariableDeclarationStatement(n, env)
	case *ast.IfStatement:
		return evaluateIfStatement(n, env)
	case *ast.WhileStatement:
		return evaluateWhileStatement(n, env)
	case *ast.ForStatement:
		return evaluateForStatement(n, env)
	case *ast.FunctionDeclaration:
		return evaluateFunctionDeclaration(n, env)
	case *ast.ReturnStatement:
		return evaluateReturnStatement(n, env)

	default:
		litter.Dump(fmt.Sprintf("Unsupported Statement node type %V\n", node))
		litter.Dump(node)
		os.Exit(65)
	}

	return nil
}

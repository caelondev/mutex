package runtime

import (
	"github.com/caelondev/mutex/src/frontend/ast"
)

func evaluateBlockStatement(block *ast.BlockStatement, env Environment) RuntimeValue {
	blockEnv := NewEnvironment(env)
	var lastEvaluated RuntimeValue = NIL()

	for _, statement := range block.Body {
		lastEvaluated = EvaluateStatement(statement, blockEnv)
	}

	return lastEvaluated
}

func evaluateVariableDeclarationStatement(stmt *ast.VariableDeclarationStatement, env Environment) RuntimeValue {
	var value RuntimeValue

	if stmt.Value != nil {
		value = EvaluateExpression(stmt.Value, env)
	} else {
		value = NIL()
	}

	return env.DeclareVariable(stmt.Identifier, value, !stmt.IsMutable)
}

func evaluateIfStatement(stmt *ast.IfStatement, env Environment) RuntimeValue {
	condition := EvaluateExpression(stmt.Condition, env)

	if isTruthy(condition) {
		return EvaluateStatement(stmt.Consequent, env)
	} else if stmt.Alternate != nil {
		return EvaluateStatement(stmt.Alternate, env)
	}

	return NIL()
}

func evaluateWhileStatement(stmt *ast.WhileStatement, env Environment) RuntimeValue {
	for {
		condition := EvaluateExpression(stmt.Condition, env)

		if !isTruthy(condition) {
			break
		}

		EvaluateStatement(stmt.Body, env)
	}

	return NIL()
}

func evaluateForStatement(stmt *ast.ForStatement, env Environment) RuntimeValue {
	loopEnv := NewEnvironment(env)

	EvaluateStatement(stmt.Initializer, loopEnv)

	for {
		condition := EvaluateExpression(stmt.Condition, loopEnv)

		if !isTruthy(condition) {
			break
		}

		EvaluateStatement(stmt.Body, loopEnv)
		EvaluateExpression(stmt.Increment, loopEnv)
	}

	return NIL()
}

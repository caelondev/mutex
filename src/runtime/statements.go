package runtime

import (
	"github.com/caelondev/mutex/src/frontend/ast"
)

func evaluateBlockStatement(block *ast.BlockStatement, env Environment) RuntimeValue {
	blockEnv := NewEnvironment(env)
	var lastEvaluated RuntimeValue = NIL()

	for _, statement := range block.Body {
		lastEvaluated = EvaluateStatement(statement, blockEnv)
		
		// If we hit a return statement, bubble it up immediately
		if _, isReturn := lastEvaluated.(*ReturnValue); isReturn {
			return lastEvaluated
		}
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

func evaluateFunctionDeclaration(stmt *ast.FunctionDeclaration, env Environment) RuntimeValue {
	functionValue := FUNCTION(stmt.Name, stmt.Parameters, stmt.Body, env)

	env.DeclareVariable(stmt.Name, functionValue, true)

	return NIL()
}

func evaluateReturnStatement(stmt *ast.ReturnStatement, env Environment) RuntimeValue {
	var value RuntimeValue

	if stmt.Value != nil {
		value = EvaluateExpression(stmt.Value, env)
	} else {
		value = NIL()
	}

	return &ReturnValue{Value: value}
}

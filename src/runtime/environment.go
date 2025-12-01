package runtime

import (
	"fmt"
	"slices"

	"github.com/caelondev/mutex/src/errors"
)

type Environment interface {
	Environment()
	DeclareVariable(variableName string, value RuntimeValue, isConstant bool) RuntimeValue
	AssignVariable(variableName string, value RuntimeValue) RuntimeValue
	ResolveVariable(variableName string) Environment
	GetVariable(variableName string) RuntimeValue
	LookupVariable(variableName string) RuntimeValue
}

type EnvironmentStruct struct {
	parent            Environment
	variables         map[string]RuntimeValue
	constantVariables []string
}

func (e *EnvironmentStruct) Environment() {}

func NewEnvironment(parentEnv Environment) *EnvironmentStruct {
	variables := map[string]RuntimeValue{}
	env := &EnvironmentStruct{
		variables: variables,
		parent:    parentEnv,
	}

	if parentEnv == nil { // This means this is the global environment
		declareGlobalVariables(env)
	}
	return env
}

func declareGlobalVariables(env Environment) {
	env.DeclareVariable("nil", NIL(), true)
	env.DeclareVariable("true", BOOLEAN(true), true)
	env.DeclareVariable("false", BOOLEAN(false), true)

	// Native functions ---
	env.DeclareVariable("echo", NATIVE_FUNCTION("echo", NATIVE_ECHO_FUNCTION), true)
	env.DeclareVariable("typeof", NATIVE_FUNCTION("typeof", NATIVE_TYPEOF_FUNCTION), true)
	env.DeclareVariable("push", NATIVE_FUNCTION("push", NATIVE_PUSH_FUNCTION), true)
	env.DeclareVariable("pop", NATIVE_FUNCTION("pop", NATIVE_POP_FUNCTION), true)
	env.DeclareVariable("shift", NATIVE_FUNCTION("shift", NATIVE_SHIFT_FUNCTION), true)
	env.DeclareVariable("unshift", NATIVE_FUNCTION("unshift", NATIVE_UNSHIFT_FUNCTION), true)

	env.DeclareVariable("string", NATIVE_FUNCTION("string", NATIVE_STRING_FUNCTION), true)
	env.DeclareVariable("int", NATIVE_FUNCTION("int", NATIVE_INT_FUNCTION), true)
	env.DeclareVariable("float", NATIVE_FUNCTION("float", NATIVE_FLOAT_FUNCTION), true)
	env.DeclareVariable("bool", NATIVE_FUNCTION("bool", NATIVE_BOOL_FUNCTION), true)

}

func (e *EnvironmentStruct) DeclareVariable(variableName string, value RuntimeValue, isConstant bool) RuntimeValue {
	if _, exists := e.variables[variableName]; exists {
		errors.ReportInterpreter(fmt.Sprintf("Cannot declare variable \"%s\" as it is already defined", variableName), 65)
	}

	if isConstant {
		e.constantVariables = append(e.constantVariables, variableName)
	}

	e.variables[variableName] = value
	return NIL()
}

func (e *EnvironmentStruct) AssignVariable(variableName string, value RuntimeValue) RuntimeValue {
	env := e.ResolveVariable(variableName)

	if envStruct, ok := env.(*EnvironmentStruct); ok {
		if slices.Contains(envStruct.constantVariables, variableName) {
			errors.ReportInterpreter(fmt.Sprintf("Cannot re-assign constant variable \"%s\"", variableName), 65)
		}
		
		envStruct.variables[variableName] = value
		return NIL()
	}

	errors.ReportInterpreter(fmt.Sprintf("Cannot re-assign variable \"%s\" as it does not exist in the current scope", variableName), 65)
	return NIL()
}

func (e *EnvironmentStruct) ResolveVariable(variableName string) Environment {
	if _, exists := e.variables[variableName]; exists {
		return e
	}

	if e.parent == nil {
		errors.ReportInterpreter(fmt.Sprintf("Cannot resolve variable \"%s\" as it does not exist in the current/outer scopes", variableName), 65)
	}

	return e.parent.ResolveVariable(variableName)
}

func (e *EnvironmentStruct) GetVariable(variableName string) RuntimeValue {
	return e.variables[variableName]
}

func (e *EnvironmentStruct) LookupVariable(variableName string) RuntimeValue {
	env := e.ResolveVariable(variableName)
	return env.GetVariable(variableName)
}

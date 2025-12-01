package runtime

import (
	"fmt"
	"strings"

	"github.com/caelondev/mutex/src/frontend/ast"
)

type ValueTypes string

const (
	BOOLEAN_VALUE ValueTypes = "boolean"
	NIL_VALUE    ValueTypes = "nil"
	NUMBER_VALUE ValueTypes = "number"
	STRING_VALUE ValueTypes = "string"
	ARRAY_VALUE ValueTypes = "array"
	FUNCTION_VALUE        ValueTypes = "function"
	NATIVE_FUNCTION_VALUE ValueTypes = "native_function"
)

type RuntimeValue interface {
	Type() ValueTypes
	String() string
}

type NilValue struct{}

func (n *NilValue) Type() ValueTypes {
	return NIL_VALUE
}

func (n *NilValue) String() string {
	return "nil"
}

type StringValue struct {
	Value string
}

func (n *StringValue) Type() ValueTypes {
	return STRING_VALUE
}

func (n *StringValue) String() string {
	return fmt.Sprintf("\"%v\"", n.Value)
}

type NumberValue struct {
	Value float64
}

func (n *NumberValue) Type() ValueTypes {
	return NUMBER_VALUE
}

func (n *NumberValue) String() string {
	return fmt.Sprintf("%v", n.Value)
}

type BooleanValue struct {
	Value bool
}

func (n *BooleanValue) Type() ValueTypes {
	return BOOLEAN_VALUE
}

func (n *BooleanValue) String() string {
	return fmt.Sprintf("%v", n.Value)
}

type FunctionValue struct {
	Name       string
	Parameters []string
	Body       ast.Statement
	Closure    Environment
}

func (f *FunctionValue) Type() ValueTypes {
	return FUNCTION_VALUE
}

func (f *FunctionValue) String() string {
	return fmt.Sprintf("[ ...function '%s'... ]", f.Name)
}


type NativeFunctionValue struct {
	Name string
	Call func(args []RuntimeValue, env Environment) RuntimeValue
}

func (n *NativeFunctionValue) Type() ValueTypes {
	return NATIVE_FUNCTION_VALUE
}

func (n *NativeFunctionValue) String() string {
	return fmt.Sprintf("[ ...native function '%s'... ]", n.Name)
}


type ReturnValue struct {
	Value RuntimeValue
}

func (r *ReturnValue) Type() ValueTypes {
	return "return"
}

func (r *ReturnValue) String() string {
	return r.Value.String()
}

type ArrayValue struct {
	Elements []RuntimeValue
}

func (a *ArrayValue) Type() ValueTypes {
	return ARRAY_VALUE
}

func (a *ArrayValue) String() string {
	if len(a.Elements) == 0 {
		return "[]"
	}

	var elements []string
	for _, elem := range a.Elements {
		elements = append(elements, fmt.Sprintf("%v", elem))
	}
	return "[" + strings.Join(elements, ", ") + "]"
}

func NIL() *NilValue {
	return &NilValue{}
}

func BOOLEAN(value bool) *BooleanValue {
	return &BooleanValue{ Value: value }
}

func ARRAY(elements []RuntimeValue) *ArrayValue {
	return &ArrayValue{Elements: elements}
}

func NATIVE_FUNCTION(name string, call func([]RuntimeValue, Environment) RuntimeValue) *NativeFunctionValue {
	return &NativeFunctionValue{
		Name: name,
		Call: call,
	}
}

func FUNCTION(name string, params []string, body ast.Statement, closure Environment) *FunctionValue {
	return &FunctionValue{
		Name:       name,
		Parameters: params,
		Body:       body,
		Closure:    closure,
	}
}

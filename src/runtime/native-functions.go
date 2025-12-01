package runtime

import (
	"fmt"
	"strconv"

	"github.com/caelondev/mutex/src/errors"
)

func NATIVE_ECHO_FUNCTION(args []RuntimeValue, env Environment) RuntimeValue {
	for i, arg := range args {
		fmt.Print(arg.String())
		if i < len(args)-1 {
			fmt.Print(" ")
		}
	}

	fmt.Println()
	return NIL()
}

func NATIVE_TYPEOF_FUNCTION(args []RuntimeValue, env Environment) RuntimeValue {
	if len(args) != 1 {
		errors.ReportInterpreter(fmt.Sprintf("typeof() expects 1 argument but got %d instead...", len(args)), 65)
	}

	return &StringValue{ Value: string(args[0].Type()) }
}

func NATIVE_PUSH_FUNCTION(args []RuntimeValue, env Environment) RuntimeValue {
	if len(args) < 2 {
		errors.ReportInterpreter("push() expects at least 2 arguments (array, value)", 65)
		return NIL()
	}

	arrayValue, ok := args[0].(*ArrayValue)
	if !ok {
		errors.ReportInterpreter(fmt.Sprintf("push() expects an array as first argument, got '%s'", args[0].Type()), 65)
		return NIL()
	}

	// Append all remaining arguments to the array
	for i := 1; i < len(args); i++ {
		arrayValue.Elements = append(arrayValue.Elements, args[i])
	}

	return NIL()
}

func NATIVE_POP_FUNCTION(args []RuntimeValue, env Environment) RuntimeValue {
	if len(args) != 1 {
		errors.ReportInterpreter("pop() expects exactly 1 argument (array)", 65)
		return NIL()
	}

	arrayValue, ok := args[0].(*ArrayValue)
	if !ok {
		errors.ReportInterpreter(fmt.Sprintf("pop() expects an array, got '%s'", args[0].Type()), 65)
		return NIL()
	}

	if len(arrayValue.Elements) == 0 {
		errors.ReportInterpreter("pop() called on empty array", 65)
		return NIL()
	}

	// Get last element
	lastElement := arrayValue.Elements[len(arrayValue.Elements)-1]
	
	// Remove last element
	arrayValue.Elements = arrayValue.Elements[:len(arrayValue.Elements)-1]

	return lastElement
}

func NATIVE_SHIFT_FUNCTION(args []RuntimeValue, env Environment) RuntimeValue {
	if len(args) != 1 {
		errors.ReportInterpreter("shift() expects exactly 1 argument (array)", 65)
		return NIL()
	}

	arrayValue, ok := args[0].(*ArrayValue)
	if !ok {
		errors.ReportInterpreter(fmt.Sprintf("shift() expects an array, got '%s'", args[0].Type()), 65)
		return NIL()
	}

	if len(arrayValue.Elements) == 0 {
		errors.ReportInterpreter("shift() called on empty array", 65)
		return NIL()
	}

	// Get first element
	firstElement := arrayValue.Elements[0]
	
	// Remove first element
	arrayValue.Elements = arrayValue.Elements[1:]

	return firstElement
}

func NATIVE_UNSHIFT_FUNCTION(args []RuntimeValue, env Environment) RuntimeValue {
	if len(args) < 2 {
		errors.ReportInterpreter("unshift() expects at least 2 arguments (array, value)", 65)
		return NIL()
	}

	arrayValue, ok := args[0].(*ArrayValue)
	if !ok {
		errors.ReportInterpreter(fmt.Sprintf("unshift() expects an array as first argument, got '%s'", args[0].Type()), 65)
		return NIL()
	}

	// Prepend all values (in order) to the beginning
	newElements := make([]RuntimeValue, 0, len(arrayValue.Elements)+len(args)-1)
	
	for i := 1; i < len(args); i++ {
		newElements = append(newElements, args[i])
	}
	
	arrayValue.Elements = append(newElements, arrayValue.Elements...)

	return NIL()
}

func NATIVE_STRING_FUNCTION(args []RuntimeValue, env Environment) RuntimeValue {
	if len(args) != 1 {
		errors.ReportInterpreter("string() expects exactly 1 argument", 65)
		return NIL()
	}

	// Convert any value to string
	return &StringValue{Value: args[0].String()}
}

func NATIVE_INT_FUNCTION(args []RuntimeValue, env Environment) RuntimeValue {
	if len(args) != 1 {
		errors.ReportInterpreter("int() expects exactly 1 argument", 65)
		return NIL()
	}

	switch v := args[0].(type) {
	case *NumberValue:
		// Truncate to integer
		return &NumberValue{Value: float64(int(v.Value))}
	case *StringValue:
		// Try to parse string as number
		parsed, err := strconv.ParseFloat(v.Value, 64)
		if err != nil {
			errors.ReportInterpreter(fmt.Sprintf("Cannot convert string '%s' to int", v.Value), 65)
			return NIL()
		}
		return &NumberValue{Value: float64(int(parsed))}
	case *BooleanValue:
		if v.Value {
			return &NumberValue{Value: 1}
		}
		return &NumberValue{Value: 0}
	default:
		errors.ReportInterpreter(fmt.Sprintf("Cannot convert type '%s' to int", v.Type()), 65)
		return NIL()
	}
}

func NATIVE_FLOAT_FUNCTION(args []RuntimeValue, env Environment) RuntimeValue {
	if len(args) != 1 {
		errors.ReportInterpreter("float() expects exactly 1 argument", 65)
		return NIL()
	}

	switch v := args[0].(type) {
	case *NumberValue:
		// Already a float
		return v
	case *StringValue:
		// Try to parse string as number
		parsed, err := strconv.ParseFloat(v.Value, 64)
		if err != nil {
			errors.ReportInterpreter(fmt.Sprintf("Cannot convert string '%s' to float", v.Value), 65)
			return NIL()
		}
		return &NumberValue{Value: parsed}
	case *BooleanValue:
		if v.Value {
			return &NumberValue{Value: 1.0}
		}
		return &NumberValue{Value: 0.0}
	default:
		errors.ReportInterpreter(fmt.Sprintf("Cannot convert type '%s' to float", v.Type()), 65)
		return NIL()
	}
}

func NATIVE_BOOL_FUNCTION(args []RuntimeValue, env Environment) RuntimeValue {
	if len(args) != 1 {
		errors.ReportInterpreter("bool() expects exactly 1 argument", 65)
		return NIL()
	}

	// Use the existing isTruthy function logic
	return BOOLEAN(isTruthy(args[0]))
}

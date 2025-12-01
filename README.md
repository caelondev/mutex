# Mutex

Mutex is a dynamically-typed scripting language implemented in Go. It features a tree-walk interpreter with support for mutable and immutable variables, control flow constructs, functions with closures, arrays, type conversions, and lexical scoping.

## Installation

### Prerequisites

- Go 1.25.3 or higher

### Installing

```bash
go install github.com/caelondev/mutex@latest
```

### Running

```bash
mutex              # Start interactive REPL
mutex <filepath>   # Execute a Mutex source file
```

## Language Reference

### Types

Mutex supports the following data types:

```mutex
// Number - IEEE 754 double-precision floating-point
42
3.14159
-69
-67.41

// String - UTF-8 encoded text
"Hello, World!"
'Single quotes work too'
`Multi-line
strings work
like this`

// Boolean - logical values
true
false

// Nil - represents absence of value
nil

// Arrays - ordered collections of values
["Mutex", 0, 1, "Hello, World", true, 3.14159]

// Functions - first-class callable objects
fn greet(name) {
    echo("Hello, " + name);
}
```

### Variable Declarations

Variables in Mutex can be declared as either mutable or immutable:

**Mutable variables** can be reassigned after declaration:

```mutex
var mut counter = 0;
counter = counter + 1;

// Compound assignment operators
counter += 1;
counter -= 1;
counter *= 2;
counter /= 2;
counter %= 3;
```

**Immutable variables** cannot be reassigned:

```mutex
var imm max_retries = 3;
// max_retries = 5;  // Error: cannot reassign immutable variable
```

### Arrays

**Mutex arrays** can be instantiated with the syntax:

```mutex
var mut numbers = [1, 2, 3, 4, 5];
var mut mixed = ["text", 42, true, nil];
var mut empty = [];
```

Arrays are reference types and cannot be directly compared:

```mutex
[] == []  // false (different references)
```

#### Indexing

Access array elements using zero-based indexing:

```mutex
var mut fruits = ["apple", "orange", "banana"];

fruits[0];  // "apple"
fruits[1];  // "orange"
fruits[2];  // "banana"
```

Nested arrays are supported:

```mutex
var mut matrix = [[1, 2], [3, 4], [5, 6]];

matrix[0][1];  // 2
matrix[2][0];  // 5
```

#### Re-assigning

Modify array elements by index:

```mutex
var mut items = ["foo", "bar"];
items[0] = "baz";  // ["baz", "bar"]
```

#### Array Methods

**push(array, ...values)** - Add elements to the end:

```mutex
var mut arr = [1, 2];
push(arr, 3);        // [1, 2, 3]
push(arr, 4, 5, 6);  // [1, 2, 3, 4, 5, 6]
```

**pop(array)** - Remove and return the last element:

```mutex
var mut arr = [1, 2, 3];
var mut last = pop(arr);  // returns 3, arr is now [1, 2]
```

**shift(array)** - Remove and return the first element:

```mutex
var mut arr = [1, 2, 3];
var mut first = shift(arr);  // returns 1, arr is now [2, 3]
```

**unshift(array, ...values)** - Add elements to the beginning:

```mutex
var mut arr = [3, 4];
unshift(arr, 2);     // [2, 3, 4]
unshift(arr, 0, 1);  // [0, 1, 2, 3, 4]
```

### Functions

Functions are first-class values with lexical closures:

```mutex
fn add(a, b) {
    return a + b;
}

var mut result = add(5, 3);  // 8
```

#### Closures

Functions capture their surrounding scope:

```mutex
fn makeCounter() {
    var mut count = 0;

    fn increment() {
        count += 1;
        return count;
    }

    return increment;
}

var mut counter = makeCounter();
counter();  // 1
counter();  // 2
counter();  // 3
```

#### Return Statements

Functions can return values explicitly or implicitly return `nil`:

```mutex
fn square(x) {
    return x * x;
}

fn noReturn() {
    echo("This returns nil");
}
```

#### Trailing Commas

Function calls support trailing commas:

```mutex
echo("a", "b", "c",);  // Valid
```

### Built-in Functions

#### echo(...values)

Prints values to stdout separated by spaces:

```mutex
echo("Hello, World!");
echo(42, "is the answer");
echo(1, 2, 3, 4, 5);
```

#### typeof(value)

Returns the type of a value as a string:

```mutex
typeof(42);           // "number"
typeof("hello");      // "string"
typeof(true);         // "boolean"
typeof(nil);          // "nil"
typeof([1, 2, 3]);    // "array"
typeof(add);          // "function"
```

#### Type Conversion Functions

**string(value)** - Convert any value to a string:

```mutex
string(42);        // "42"
string(true);      // "true"
string([1, 2]);    // "[1, 2]"
string(3.14);      // "3.14"
```

**int(value)** - Convert to integer (truncates decimals):

```mutex
int(3.14);         // 3
int("42");         // 42
int(true);         // 1
int(false);        // 0
```

**float(value)** - Convert to floating-point number:

```mutex
float("3.14");     // 3.14
float(42);         // 42.0
float(true);       // 1.0
float(false);      // 0.0
```

**bool(value)** - Convert to boolean (using truthiness):

```mutex
bool(1);           // true
bool(0);           // false
bool("hello");     // true
bool("");          // false
bool(nil);         // false
bool([]);          // true
```

### Expressions

#### Arithmetic Operators

```mutex
10 + 5    // Addition: 15
10 - 5    // Subtraction: 5
10 * 5    // Multiplication: 50
10 / 5    // Division: 2
10 % 3    // Modulo: 1
```

#### Compound Assignment Operators

```mutex
x += 5    // Shorthand for x = x + 5
x -= 5    // Shorthand for x = x - 5
x *= 5    // Shorthand for x = x * 5
x /= 5    // Shorthand for x = x / 5
x %= 3    // Shorthand for x = x % 3
```

#### Increment/Decrement Operators

Postfix operators only:

```mutex
x++;      // Post-increment: returns x, then adds 1
x--;      // Post-decrement: returns x, then subtracts 1
```

Example:

```mutex
var mut i = 5;
echo(i++);  // prints 5, i becomes 6
echo(i);    // prints 6
```

#### Comparison Operators

```mutex
5 < 10     // Less than: true
5 <= 10    // Less than or equal: true
5 > 10     // Greater than: false
5 >= 10    // Greater than or equal: false
5 == 10    // Equal to: false
5 != 10    // Not equal to: true
```

#### Logical Operators

Mutex uses keyword-based logical operators with short-circuit evaluation:

```mutex
true and true         // true (both must be true)
false or true         // true (at least one must be true)
not true              // false (negation)
not (true and true)   // false

// Short-circuit evaluation
false and echo("This won't print")  // false, echo never runs
true or echo("This won't print")    // true, echo never runs
```

#### String Concatenation

```mutex
"Hello" + " " + "World"  // "Hello World"

// Mix with type conversion
string(5) + " items"     // "5 items"
"Value: " + string(3.14) // "Value: 3.14"
```

### Control Flow

#### Conditional Statements

Execute code based on boolean conditions:

```mutex
if (temperature > 30) {
    echo("Hot weather");
} else if (temperature > 20) {
    echo("Mild weather");
} else {
    echo("Cold weather");
}
```

Parentheses around conditions are optional:

```mutex
if temperature > 30 {
    echo("Also valid");
}
```

#### While Loops

Execute code repeatedly while a condition holds true:

```mutex
var mut i = 0;
while (i < 5) {
    echo(i);
    i += 1;
}
```

#### For Loops

Execute code with explicit initialization, condition, and increment:

```mutex
for (var mut i = 0; i < 10; i++) {
    echo(i);
}
// i is not accessible here (scoped to loop)
```

### Scoping Rules

Mutex uses lexical scoping:

- Variables declared outside blocks are accessible within nested blocks
- Variables declared within blocks are not accessible outside those blocks
- For loop initializers are scoped to the loop body
- Functions capture their closure environment

```mutex
var mut x = 10;

if (x > 5) {
    var mut y = 20;  // y is scoped to this block
    x += y;          // x is accessible from outer scope
}

// x is 30
// y is not accessible here
```

### Truthiness

Values are considered "truthy" or "falsy" in boolean contexts:

**Falsy values:**

- `nil`
- `false`
- `0` (number zero)
- `""` (empty string)

**Truthy values:**

- Everything else, including:
  - Non-zero numbers
  - Non-empty strings
  - All arrays (even empty ones)
  - All functions

## Examples

### Sum of Natural Numbers

```mutex
var mut sum = 0;
for (var mut i = 1; i <= 100; i++) {
    sum += i;
}
echo(sum);  // 5050
```

### Factorial Calculation

```mutex
fn factorial(n) {
    if (n <= 1) {
        return 1;
    }
    return n * factorial(n - 1);
}

echo(factorial(5));  // 120
```

### FizzBuzz

```mutex
for (var mut i = 1; i <= 15; i++) {
    var mut output = "";

    if (i % 3 == 0) {
        output = "Fizz";
    }

    if (i % 5 == 0) {
        output = output + "Buzz";
    }

    if (output == "") {
        output = string(i);
    }

    echo(output);
}
```

### Array Manipulation

```mutex
var mut numbers = [1, 2, 3, 4, 5];

// Calculate sum
var mut sum = 0;
for (var mut i = 0; i < 5; i++) {
    sum += numbers[i];
}
echo("Sum:", sum);  // Sum: 15

// Double all values
for (var mut i = 0; i < 5; i++) {
    numbers[i] = numbers[i] * 2;
}
echo(numbers);  // [2, 4, 6, 8, 10]
```

### Using Array Methods

```mutex
var mut stack = [];

push(stack, 1);
push(stack, 2);
push(stack, 3);
echo(stack);  // [1, 2, 3]

var mut top = pop(stack);
echo("Popped:", top);  // Popped: 3
echo(stack);           // [1, 2]

var mut queue = [];
push(queue, "first");
push(queue, "second");
push(queue, "third");

var mut item = shift(queue);
echo("Processed:", item);  // Processed: first
echo(queue);               // ["second", "third"]
```

### Higher-Order Functions

```mutex
fn map(arr, func) {
    var mut result = [];
    for (var mut i = 0; i < 5; i++) {
        push(result, func(arr[i]));
    }
    return result;
}

fn double(x) {
    return x * 2;
}

var mut numbers = [1, 2, 3, 4, 5];
var mut doubled = map(numbers, double);
echo(doubled);  // [2, 4, 6, 8, 10]
```

### Type Conversion Examples

```mutex
// Building formatted strings
var mut age = 25;
var mut message = "You are " + string(age) + " years old";
echo(message);  // "You are 25 years old"

// Parsing user input (simulated)
var mut input = "42";
var mut number = int(input);
echo(number * 2);  // 84

// Converting booleans
var mut hasAccess = bool(1);
if (hasAccess) {
    echo("Access granted");
}
```

## Implementation

Mutex is implemented as a tree-walk interpreter consisting of:

- **Lexer** - Tokenizes source code into lexical tokens
- **Parser** - Constructs an abstract syntax tree (AST) using Pratt parsing
- **Evaluator** - Traverses and executes the AST with environment-based variable storage

The interpreter supports:

- Proper lexical scoping through environment chaining
- First-class functions with closure support
- Dynamic typing with runtime type checking
- Short-circuit evaluation for logical operators
- Reference semantics for arrays (mutation in-place)

## Author

Built by [caelondev](https://github.com/caelondev)

## License

MIT

package parser

import (
	"github.com/caelondev/mutex/src/frontend/ast"
	"github.com/caelondev/mutex/src/frontend/lexer"
)

func parseStatement(p *parser) ast.Statement {
	if p.isEOF() {
		return nil
	}
	statementFunction, exists := statementLU[p.currentTokenType()]

	if exists {
		return statementFunction(p)
	}

	expression := parseExpression(p, DEFAULT_BP)

	p.expect(lexer.SEMICOLON)

	return &ast.ExpressionStatement{
		Expression: expression,
	}
}

func parseVariableDeclaration(p *parser) ast.Statement {
	//  SYNTAX ---
	//
	//  var (mut | imm) variableName = value
	//

	var identifier string
	var value ast.Expression
	var isMutable bool

	p.advance() // eat var keyword ---

	// check (imm/mut)
	mutabilityType := p.expect(lexer.MUTABLE, lexer.IMMUTABLE).TokenType

	isMutable = mutabilityType == lexer.MUTABLE

	identifier = p.expect(lexer.IDENTIFIER).Lexeme

	if p.currentTokenType() != lexer.SEMICOLON {
		p.expect(lexer.ASSIGNMENT, lexer.SEMICOLON)
		// NOTE: EXPECTING A SEMICOLON HERE IS USELESS... I JUST USED IT FOR ERROR MESSAGE ---

		value = parseExpression(p, DEFAULT_BP)
	}

	p.expect(lexer.SEMICOLON)

	return &ast.VariableDeclarationStatement{
		Identifier: identifier,
		IsMutable:  isMutable,
		Value:      value,
	}
}

func parseIfStatement(p *parser) ast.Statement {
	//  SYNTAX ---
	//
	//  if (condition) { ... } else { ... }
	//  if (condition) { ... } else if (condition) { ... }
	//

	p.advance() // eat 'if' keyword

	// Parse condition (can be with or without parentheses)
	p.ignore(lexer.LEFT_PARENTHESIS)

	condition := parseExpression(p, DEFAULT_BP)

	p.ignore(lexer.RIGHT_PARENTHESIS) // eat ')'

	// Parse consequent block
	p.expect(lexer.LEFT_BRACE)
	consequent := parseBlock(p)

	// Parse optional else/else if
	var alternate ast.Statement
	if p.currentTokenType() == lexer.ELSE {
		p.advance() // eat 'else'

		// Check if it's 'else if' or just 'else'
		if p.currentTokenType() == lexer.IF {
			alternate = parseIfStatement(p) // recursive for else if
		} else {
			p.expect(lexer.LEFT_BRACE)
			alternate = parseBlock(p)
		}
	}

	return &ast.IfStatement{
		Condition:  condition,
		Consequent: consequent,
		Alternate:  alternate,
	}
}

func parseBlock(p *parser) ast.Statement {
	// Assumes LEFT_BRACE already consumed
	var body []ast.Statement

	for !p.isEOF() && p.currentTokenType() != lexer.RIGHT_BRACE {
		statement := parseStatement(p)
		if statement != nil {
			body = append(body, statement)
		}
	}

	p.expect(lexer.RIGHT_BRACE)

	return &ast.BlockStatement{
		Body: body,
	}
}

func parseWhileStatement (p *parser) ast.Statement {
	// SYNTAX ---
	//
	// while(condition) { ... } ---
	//
	
	p.advance() // eat while token ---
	p.ignore(lexer.LEFT_PARENTHESIS)

	condition := parseExpression(p, DEFAULT_BP)

	p.ignore(lexer.RIGHT_PARENTHESIS)
	p.expect(lexer.LEFT_BRACE)

	body := parseBlock(p)
	
	return &ast.WhileStatement{
		Condition: condition,
		Body: body,
	}

}

func parseForStatement(p *parser) ast.Statement {
	// SYNTAX --- 
	//
	// for (initializer; condition; increment) { ... } ---
	//

	p.advance() // Eat `for` token

	p.expect(lexer.LEFT_PARENTHESIS)

	initializer := parseVariableDeclaration(p) // Already consumes semicolon

	condition := parseExpression(p, DEFAULT_BP)
	p.expect(lexer.SEMICOLON)

	increment := parseExpression(p, DEFAULT_BP)

	p.expect(lexer.RIGHT_PARENTHESIS)

	// Parse body
	p.expect(lexer.LEFT_BRACE)
	body := parseBlock(p)

	return &ast.ForStatement{
		Initializer: initializer,
		Condition:   condition,
		Increment:   increment,
		Body:        body,
	}
}

func parseFunctionDeclaration(p *parser) ast.Statement {
	// SYNTAX ---
	//
	// fn name(param1, param2, ...) { ... }
	//
	
	var parameters []string

	p.advance() // Eat 'fn' ---
	name := p.expect(lexer.IDENTIFIER).Lexeme

	// Parse Parameters ---
	p.expect(lexer.LEFT_PARENTHESIS)

	if p.currentTokenType() != lexer.RIGHT_PARENTHESIS {
		parameters = append(parameters, p.expect(lexer.IDENTIFIER).Lexeme)
		for p.currentTokenType() == lexer.COMMA {
			p.advance() // eat ',' ---
			parameters = append(parameters, p.expect(lexer.IDENTIFIER).Lexeme)
		}
	}

	p.expect(lexer.RIGHT_PARENTHESIS)

	// Parse Function Body ---
	p.expect(lexer.LEFT_BRACE)
	body := parseBlock(p)

	return &ast.FunctionDeclaration{
		Name: name,
		Parameters: parameters,
		Body: body,
	}
}

func parseReturnStatement(p *parser) ast.Statement {
	// SYNTAX ---
	//
	// return;
	// return value;
	// return x+y;

	p.advance() // Eat 'return' keyword ---

	var value ast.Expression

	// Check return value, not just 'return' itself
	if p.currentTokenType() != lexer.SEMICOLON && !p.isEOF() {
		value = parseExpression(p, DEFAULT_BP)
	}

	p.expect(lexer.SEMICOLON)

	return &ast.ReturnStatement{
		Value: value,
	}
}

func parseCallExpression(p *parser, left ast.Expression, bp BindingPower) ast.Expression {
	// SYNTAX ---
	//
	// functionName(arg1, arg2, arg3)
	// functionName(arg1, arg2,)  // trailing comma allowed
	
	p.advance() // eat '(' ---
	
	var args []ast.Expression
	
	// Parse arguments
	if p.currentTokenType() != lexer.RIGHT_PARENTHESIS {
		args = append(args, parseExpression(p, DEFAULT_BP))
		
		for p.currentTokenType() == lexer.COMMA {
			p.advance() // eat ','
			
			// Check for trailing comma before RIGHT_PARENTHESIS
			if p.currentTokenType() == lexer.RIGHT_PARENTHESIS {
				break
			}
			
			args = append(args, parseExpression(p, DEFAULT_BP))
		}
	}
	
	p.expect(lexer.RIGHT_PARENTHESIS)
	
	return &ast.CallExpression{
		Callee:    left,
		Arguments: args,
	}
}

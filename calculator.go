package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

// TokenType represents the type of a lexer token.
type TokenType int

const (
	TokenNumber TokenType = iota
	TokenPlus
	TokenMinus
	TokenMul
	TokenDiv
	TokenLParen
	TokenRParen
)

// Token represents a single token produced by the lexer.
type Token struct {
	Type  TokenType
	Value string
}

// tokenize splits an expression string into a slice of tokens.
func tokenize(expr string) ([]Token, error) {
	var tokens []Token
	expr = strings.TrimSpace(expr)
	i := 0

	for i < len(expr) {
		ch := rune(expr[i])

		if unicode.IsSpace(ch) {
			i++
			continue
		}

		switch ch {
		case '+':
			tokens = append(tokens, Token{Type: TokenPlus, Value: "+"})
			i++
		case '-':
			// Handle unary minus: at the start, after '(' or after an operator
			if len(tokens) == 0 || tokens[len(tokens)-1].Type == TokenLParen ||
				tokens[len(tokens)-1].Type == TokenPlus ||
				tokens[len(tokens)-1].Type == TokenMinus ||
				tokens[len(tokens)-1].Type == TokenMul ||
				tokens[len(tokens)-1].Type == TokenDiv {
				// Parse as a negative number
				i++
				start := i
				for i < len(expr) && (expr[i] == '.' || (expr[i] >= '0' && expr[i] <= '9')) {
					i++
				}
				if i == start {
					return nil, fmt.Errorf("ungültiges Zeichen nach '-' an Position %d", start)
				}
				tokens = append(tokens, Token{Type: TokenNumber, Value: "-" + expr[start:i]})
			} else {
				tokens = append(tokens, Token{Type: TokenMinus, Value: "-"})
				i++
			}
		case '*':
			tokens = append(tokens, Token{Type: TokenMul, Value: "*"})
			i++
		case '/':
			tokens = append(tokens, Token{Type: TokenDiv, Value: "/"})
			i++
		case '(':
			tokens = append(tokens, Token{Type: TokenLParen, Value: "("})
			i++
		case ')':
			tokens = append(tokens, Token{Type: TokenRParen, Value: ")"})
			i++
		default:
			if ch == '.' || (ch >= '0' && ch <= '9') {
				start := i
				dotCount := 0
				for i < len(expr) && (expr[i] == '.' || (expr[i] >= '0' && expr[i] <= '9')) {
					if expr[i] == '.' {
						dotCount++
					}
					i++
				}
				if dotCount > 1 {
					return nil, fmt.Errorf("ungültige Zahl: %s", expr[start:i])
				}
				tokens = append(tokens, Token{Type: TokenNumber, Value: expr[start:i]})
			} else {
				return nil, fmt.Errorf("unbekanntes Zeichen: '%c'", ch)
			}
		}
	}
	return tokens, nil
}

// precedence returns the operator precedence.
func precedence(t TokenType) int {
	switch t {
	case TokenPlus, TokenMinus:
		return 1
	case TokenMul, TokenDiv:
		return 2
	default:
		return 0
	}
}

// isOperator returns true if the token type is an operator.
func isOperator(t TokenType) bool {
	return t == TokenPlus || t == TokenMinus || t == TokenMul || t == TokenDiv
}

// shuntingYard converts infix tokens to postfix (Reverse Polish Notation) using the
// Shunting-Yard algorithm.
func shuntingYard(tokens []Token) ([]Token, error) {
	var output []Token
	var opStack []Token

	for _, tok := range tokens {
		switch {
		case tok.Type == TokenNumber:
			output = append(output, tok)

		case isOperator(tok.Type):
			for len(opStack) > 0 {
				top := opStack[len(opStack)-1]
				if isOperator(top.Type) && precedence(top.Type) >= precedence(tok.Type) {
					output = append(output, top)
					opStack = opStack[:len(opStack)-1]
				} else {
					break
				}
			}
			opStack = append(opStack, tok)

		case tok.Type == TokenLParen:
			opStack = append(opStack, tok)

		case tok.Type == TokenRParen:
			found := false
			for len(opStack) > 0 {
				top := opStack[len(opStack)-1]
				opStack = opStack[:len(opStack)-1]
				if top.Type == TokenLParen {
					found = true
					break
				}
				output = append(output, top)
			}
			if !found {
				return nil, errors.New("fehlende öffnende Klammer")
			}
		}
	}

	for len(opStack) > 0 {
		top := opStack[len(opStack)-1]
		opStack = opStack[:len(opStack)-1]
		if top.Type == TokenLParen {
			return nil, errors.New("fehlende schließende Klammer")
		}
		output = append(output, top)
	}

	return output, nil
}

// evalPostfix evaluates a postfix (RPN) token sequence and returns the result.
func evalPostfix(tokens []Token) (float64, error) {
	var stack []float64

	for _, tok := range tokens {
		if tok.Type == TokenNumber {
			val, err := strconv.ParseFloat(tok.Value, 64)
			if err != nil {
				return 0, fmt.Errorf("ungültige Zahl: %s", tok.Value)
			}
			stack = append(stack, val)
		} else if isOperator(tok.Type) {
			if len(stack) < 2 {
				return 0, errors.New("ungültiger Ausdruck")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var result float64
			switch tok.Type {
			case TokenPlus:
				result = a + b
			case TokenMinus:
				result = a - b
			case TokenMul:
				result = a * b
			case TokenDiv:
				if b == 0 {
					return 0, errors.New("Division durch Null")
				}
				result = a / b
			}
			stack = append(stack, result)
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("ungültiger Ausdruck")
	}

	result := stack[0]
	if math.IsInf(result, 0) || math.IsNaN(result) {
		return 0, errors.New("ungültiges Ergebnis")
	}

	return result, nil
}

// Calculate evaluates a mathematical expression string and returns the result.
func Calculate(expression string) (float64, error) {
	expression = strings.TrimSpace(expression)
	if expression == "" {
		return 0, errors.New("leerer Ausdruck")
	}

	tokens, err := tokenize(expression)
	if err != nil {
		return 0, err
	}

	if len(tokens) == 0 {
		return 0, errors.New("leerer Ausdruck")
	}

	postfix, err := shuntingYard(tokens)
	if err != nil {
		return 0, err
	}

	return evalPostfix(postfix)
}


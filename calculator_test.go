package main

import (
	"math"
	"testing"
)

func TestCalculateBasicOperations(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		expected   float64
	}{
		{"addition", "2 + 3", 5},
		{"subtraction", "10 - 4", 6},
		{"multiplication", "6 * 7", 42},
		{"division", "20 / 4", 5},
		{"decimal addition", "1.5 + 2.5", 4},
		{"decimal multiplication", "2.5 * 4", 10},
		{"negative result", "3 - 10", -7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Calculate(tt.expression)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("Calculate(%q) = %v, want %v", tt.expression, result, tt.expected)
			}
		})
	}
}

func TestCalculatePrecedence(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		expected   float64
	}{
		{"mul before add", "2 + 3 * 4", 14},
		{"div before sub", "10 - 6 / 3", 8},
		{"left to right mul div", "12 / 3 * 2", 8},
		{"complex precedence", "2 + 3 * 4 - 1", 13},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Calculate(tt.expression)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("Calculate(%q) = %v, want %v", tt.expression, result, tt.expected)
			}
		})
	}
}

func TestCalculateParentheses(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		expected   float64
	}{
		{"simple parens", "(2 + 3) * 4", 20},
		{"nested parens", "((2 + 3) * (4 - 1))", 15},
		{"parens override precedence", "12.5 * (3 + 4)", 87.5},
		{"deep nesting", "((1 + 2) * (3 + (4 - 2)))", 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Calculate(tt.expression)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("Calculate(%q) = %v, want %v", tt.expression, result, tt.expected)
			}
		})
	}
}

func TestCalculateUnaryMinus(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		expected   float64
	}{
		{"negative number", "-5", -5},
		{"negative in expression", "-5 + 3", -2},
		{"negative after paren", "(-3) * 2", -6},
		{"double negative", "-3 + -2", -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Calculate(tt.expression)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("Calculate(%q) = %v, want %v", tt.expression, result, tt.expected)
			}
		})
	}
}

func TestCalculateErrors(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		errContain string
	}{
		{"division by zero", "10 / 0", "Division durch Null"},
		{"empty expression", "", "leerer Ausdruck"},
		{"whitespace only", "   ", "leerer Ausdruck"},
		{"invalid character", "2 + a", "unbekanntes Zeichen"},
		{"missing closing paren", "(2 + 3", "schließende Klammer"},
		{"missing opening paren", "2 + 3)", "öffnende Klammer"},
		{"invalid number", "1.2.3", "ungültige Zahl"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Calculate(tt.expression)
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tt.errContain)
			}
			if !contains(err.Error(), tt.errContain) {
				t.Errorf("error %q does not contain %q", err.Error(), tt.errContain)
			}
		})
	}
}

func TestCalculateEdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		expected   float64
	}{
		{"single number", "42", 42},
		{"decimal number", "3.14", 3.14},
		{"no spaces", "2+3*4", 14},
		{"extra spaces", "  2  +  3  ", 5},
		{"zero result", "5 - 5", 0},
		{"large expression", "1 + 2 + 3 + 4 + 5 + 6 + 7 + 8 + 9 + 10", 55},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Calculate(tt.expression)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("Calculate(%q) = %v, want %v", tt.expression, result, tt.expected)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}


package database

import "fmt"

type (
	// Expression is a mixed value type
	Expression struct {
		value interface{}
	}
)

// NewExpression returns a new Expression
func NewExpression(value interface{}) *Expression {
	return &Expression{
		value: value,
	}
}

// GetValue returns the value
func (ex *Expression) GetValue() interface{} {
	return ex.value
}

// ToString converts value to string
func (ex *Expression) ToString() string {
	return fmt.Sprintf("%v", ex.value)
}

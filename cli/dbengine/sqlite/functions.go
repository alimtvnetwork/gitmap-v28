package sqlite

import (
	"fmt"
	"strings"
)

// CompileFunctionCall builds a SELECT <func>(args...) query to run a database function.
func (c *Compiler) CompileFunctionCall(name string, argCount int) string {
	if argCount <= 0 {
		return fmt.Sprintf("SELECT %s();", name)
	}

	placeholders := make([]string, argCount)
	for i := 0; i < argCount; i++ {
		placeholders[i] = "?"
	}

	return fmt.Sprintf("SELECT %s(%s);", name, strings.Join(placeholders, ", "))
}

// CompileScalarFunctionExpression returns a scalar function expression string like datetime(?, ?).
func (c *Compiler) CompileScalarFunctionExpression(name string, argCount int) string {
	if argCount <= 0 {
		return fmt.Sprintf("%s()", name)
	}

	placeholders := make([]string, argCount)
	for i := 0; i < argCount; i++ {
		placeholders[i] = "?"
	}

	return fmt.Sprintf("%s(%s)", name, strings.Join(placeholders, ", "))
}

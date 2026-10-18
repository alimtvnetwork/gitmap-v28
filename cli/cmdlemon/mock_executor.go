package cmdlemon

import (
	"context"
	"errors"
)

// MockCommandExecutor provides an in-memory executor for testing.
type MockCommandExecutor struct {
	Handler func(ctx context.Context, args []string) (ExecutionResult, error)
	Calls   [][]string
}

func (m *MockCommandExecutor) Execute(ctx context.Context, raw string) (ExecutionResult, error) {
	args, err := ParseAndAuthorizeCommand(raw)
	if err != nil {
		return ExecutionResult{ExitCode: 1, Stderr: err.Error()}, err
	}
	return m.ExecuteArgs(ctx, args)
}

func (m *MockCommandExecutor) ExecuteArgs(ctx context.Context, args []string) (ExecutionResult, error) {
	validArgs, ok := ValidateTokens(args)
	if !ok {
		return ExecutionResult{ExitCode: 1, Stderr: "unauthorized"}, errors.New("unauthorized command")
	}
	m.Calls = append(m.Calls, validArgs)
	if m.Handler != nil {
		return m.Handler(ctx, validArgs)
	}
	return ExecutionResult{ExitCode: 0, Stdout: "mock success", IsSuccess: true}, nil
}

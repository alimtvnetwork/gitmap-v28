package netip

import (
	"context"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ConnectivityValidator checks ICMP reachability for gateways and internet hosts.
type ConnectivityValidator struct{}

// NewConnectivityValidator creates a new instance of ConnectivityValidator.
func NewConnectivityValidator() *ConnectivityValidator {
	return &ConnectivityValidator{}
}

// ValidateConnectivity checks ICMP reachability according to options.
func (v *ConnectivityValidator) ValidateConnectivity(ctx context.Context, opts ValidationOptions) *apperror.AppError {
	if !opts.IsValidationActive {
		return nil
	}

	normOpts := normalizeValidationOptions(opts)
	if err := validateGateway(ctx, normOpts); err != nil {
		return err
	}

	return pingHost(ctx, normOpts.TestIP, normOpts.PacketCount, normOpts.TimeoutSeconds)
}

func validateGateway(ctx context.Context, opts ValidationOptions) *apperror.AppError {
	if !opts.IsPingGateway || opts.Gateway == "" {
		return nil
	}

	return pingHost(ctx, opts.Gateway, opts.PacketCount, opts.TimeoutSeconds)
}

func normalizeValidationOptions(opts ValidationOptions) ValidationOptions {
	if opts.TestIP == "" {
		opts.TestIP = "8.8.8.8"
	}
	if opts.PacketCount <= 0 {
		opts.PacketCount = 2
	}
	if opts.TimeoutSeconds <= 0 {
		opts.TimeoutSeconds = 3
	}

	return opts
}

func pingHost(ctx context.Context, target string, count, timeoutSec int) *apperror.AppError {
	cmd := buildPingCmd(ctx, target, count, timeoutSec)
	if err := cmd.Run(); err != nil {
		return apperror.WrapWithDetails(err, "pingHost", "E_PING_FAILED", "host unreachable: "+target, "validator", apperror.ErrorTypeExecution, apperror.SeverityError, map[string]any{"target": target})
	}

	return nil
}

func buildPingCmd(ctx context.Context, target string, count, timeoutSec int) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return buildWindowsPingCmd(ctx, target, count, timeoutSec)
	}

	return buildLinuxPingCmd(ctx, target, count, timeoutSec)
}

func buildWindowsPingCmd(ctx context.Context, target string, count, timeoutSec int) *exec.Cmd {
	timeoutMs := strconv.Itoa(timeoutSec * 1000)

	return exec.CommandContext(ctx, "ping", "-n", strconv.Itoa(count), "-w", timeoutMs, target)
}

func buildLinuxPingCmd(ctx context.Context, target string, count, timeoutSec int) *exec.Cmd {
	return exec.CommandContext(ctx, "ping", "-c", strconv.Itoa(count), "-W", strconv.Itoa(timeoutSec), target)
}

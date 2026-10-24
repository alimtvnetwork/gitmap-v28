package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// SSHDiagnosisType classifies the root-cause diagnosis of SSH connectivity probes.
type SSHDiagnosisType string

const (
	DiagFirewallInboundDrop SSHDiagnosisType = "FIREWALL INBOUND DROP"
	DiagDaemonNotRunning    SSHDiagnosisType = "SSH DAEMON NOT RUNNING"
	DiagHostOffline         SSHDiagnosisType = "HOST UNREACHABLE / OFFLINE"
	DiagSSHPortOpen         SSHDiagnosisType = "SSH PORT OPEN & ACCESSIBLE"
	DiagUnknownError        SSHDiagnosisType = "NETWORK PROBE ERROR"
)

// IsValid reports whether the diagnosis matches a known diagnosis value.
func (d SSHDiagnosisType) IsValid() bool {
	switch d {
	case DiagFirewallInboundDrop, DiagDaemonNotRunning, DiagHostOffline,
		DiagSSHPortOpen, DiagUnknownError:
		return true
	default:
		return false
	}
}

// String returns the string representation of the diagnosis.
func (d SSHDiagnosisType) String() string {
	return string(d)
}

// ParseSSHDiagnosisType parses a raw string into a validated SSHDiagnosisType.
// Diagnosis values are uppercase display strings, so matching is
// case-insensitive rather than the lowercase normalization used for
// lowercase-valued enums.
func ParseSSHDiagnosisType(s string) (SSHDiagnosisType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, candidate := range AllSSHDiagnosisTypes() {
		if strings.EqualFold(trimmed, candidate.String()) {
			return candidate, nil
		}
	}

	return "", apperror.NewValidation("ParseSSHDiagnosisType", "E1000", fmt.Sprintf("invalid ssh diagnosis: %q", s)).WithContext("value", s)
}

// AllSSHDiagnosisTypes returns all valid SSHDiagnosisType variants.
func AllSSHDiagnosisTypes() []SSHDiagnosisType {
	return []SSHDiagnosisType{
		DiagFirewallInboundDrop,
		DiagDaemonNotRunning,
		DiagHostOffline,
		DiagSSHPortOpen,
		DiagUnknownError,
	}
}

// IsOpen reports whether SSH port is open and accessible.
func (d SSHDiagnosisType) IsOpen() bool {
	return d == DiagSSHPortOpen
}

// IsFirewallInboundDrop reports whether inbound traffic is firewall-dropped.
func (d SSHDiagnosisType) IsFirewallInboundDrop() bool {
	return d == DiagFirewallInboundDrop
}

// IsDaemonNotRunning reports whether the SSH daemon is not running.
func (d SSHDiagnosisType) IsDaemonNotRunning() bool {
	return d == DiagDaemonNotRunning
}

// IsHostOffline reports whether the host is unreachable or offline.
func (d SSHDiagnosisType) IsHostOffline() bool {
	return d == DiagHostOffline
}

// IsUnknownError reports whether the diagnosis is a network probe error.
func (d SSHDiagnosisType) IsUnknownError() bool {
	return d == DiagUnknownError
}

// HasDiagnosis reports whether the diagnosis is non-empty and valid.
func (d SSHDiagnosisType) HasDiagnosis() bool {
	return d.IsValid()
}

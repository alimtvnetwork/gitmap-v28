package netip

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type (
	// InterfaceInfo represents network adapter state and configuration.
	InterfaceInfo struct {
		Name       string   `json:"name"`
		IP         string   `json:"ip"`
		Netmask    string   `json:"netmask"`
		Gateway    string   `json:"gateway"`
		DNS        []string `json:"dns"`
		MAC        string   `json:"mac"`
		IsUp       bool     `json:"isUp"`
		IsLoopback bool     `json:"isLoopback"`
		IsDHCP     bool     `json:"isDHCP"`
		IsValid    bool     `json:"isValid"`
	}

	// ChangeOptions specifies mutations for IP interface configuration.
	ChangeOptions struct {
		InterfaceName string   `json:"interfaceName"`
		IP            string   `json:"ip"`
		Netmask       string   `json:"netmask"`
		Gateway       string   `json:"gateway"`
		DNS           []string `json:"dns"`
		IsDHCP        bool     `json:"isDHCP"`
		IsAutoConfirm bool     `json:"isAutoConfirm"`
		IsDryRun      bool     `json:"isDryRun"`
		IsValid       bool     `json:"isValid"`
	}

	// Snapshot records an interface configuration state for rollbacks.
	Snapshot struct {
		IPSnapshotId  int64    `json:"ipSnapshotId"`
		InterfaceName string   `json:"interfaceName"`
		IP            string   `json:"ip"`
		Netmask       string   `json:"netmask"`
		Gateway       string   `json:"gateway"`
		DNS           []string `json:"dns"`
		IsDHCP        bool     `json:"isDHCP"`
		Timestamp     int64    `json:"timestamp"`
		Notes         string   `json:"notes,omitempty"`
		Comments      string   `json:"comments,omitempty"`
		IsValid       bool     `json:"isValid"`
	}

	// RollbackResult contains the outcome of an IP rollback operation.
	RollbackResult struct {
		IsReverted    bool   `json:"isReverted"`
		InterfaceName string `json:"interfaceName"`
		RestoredIP    string `json:"restoredIP"`
		Message       string `json:"message"`
		IsValid       bool   `json:"isValid"`
	}

	// ValidationOptions defines connectivity verification settings.
	ValidationOptions struct {
		Gateway            string `json:"gateway"`
		TestIP             string `json:"testIP"`
		TimeoutSeconds     int    `json:"timeoutSeconds"`
		PacketCount        int    `json:"packetCount"`
		IsPingGateway      bool   `json:"isPingGateway"`
		IsValidationActive bool   `json:"isValidationActive"`
		IsValid            bool   `json:"isValid"`
	}

	// InterfaceResult wraps single InterfaceInfo response.
	InterfaceResult = result.Result[InterfaceInfo]

	// InterfaceSliceResult wraps multiple InterfaceInfo response.
	InterfaceSliceResult = result.ResultSlice[InterfaceInfo]

	// RollbackResultWrap wraps RollbackResult response.
	RollbackResultWrap = result.Result[RollbackResult]

	// SnapshotResult wraps single Snapshot response.
	SnapshotResult = result.Result[Snapshot]

	// SnapshotSliceResult wraps multiple Snapshot response.
	SnapshotSliceResult = result.ResultSlice[Snapshot]
)

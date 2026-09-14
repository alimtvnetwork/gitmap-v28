package netip

import (
	"context"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// Manager orchestrates network IP mutations, snapshots, and rollbacks.
type Manager struct {
	driver    Driver
	db        *store.DB
	validator *ConnectivityValidator
}

// NewManager creates an instance of Manager.
func NewManager(driver Driver, db *store.DB) *Manager {
	return &Manager{
		driver:    driver,
		db:        db,
		validator: NewConnectivityValidator(),
	}
}

// Driver returns the currently configured driver.
func (m *Manager) Driver() Driver {
	return m.driver
}

// ListInterfaces lists all network interfaces.
func (m *Manager) ListInterfaces(ctx context.Context) InterfaceSliceResult {
	if m.driver == nil {
		return result.FailSlice[InterfaceInfo](apperror.NewExecutionError("network IP driver not available on this platform"))
	}

	return m.driver.ListInterfaces(ctx)
}

// GetInterface retrieves a specific network interface.
func (m *Manager) GetInterface(ctx context.Context, name string) InterfaceResult {
	if m.driver == nil {
		return result.Fail[InterfaceInfo](apperror.NewExecutionError("network IP driver not available on this platform"))
	}

	return m.driver.GetInterface(ctx, name)
}

// ChangeIP applies new IP configuration with pre-snapshotting and validation.
func (m *Manager) ChangeIP(
	ctx context.Context,
	opts ChangeOptions,
	valOpts ValidationOptions,
) RollbackResultWrap {
	if m.driver == nil {
		return result.Fail[RollbackResult](apperror.NewExecutionError("network IP driver not available on this platform"))
	}

	if err := validateChangeOptions(opts); err != nil {
		return result.Fail[RollbackResult](err)
	}

	snap, err := m.prepareSnapshot(ctx, opts.InterfaceName)
	if err != nil {
		return result.Fail[RollbackResult](err)
	}

	return m.dispatchChange(ctx, opts, valOpts, snap)
}

func (m *Manager) dispatchChange(
	ctx context.Context,
	opts ChangeOptions,
	valOpts ValidationOptions,
	snap Snapshot,
) RollbackResultWrap {
	if opts.IsDryRun {
		return buildDryRunResult(opts.InterfaceName, snap.IP)
	}

	return m.executeChangeWithRollback(ctx, opts, valOpts, snap)
}

func validateChangeOptions(opts ChangeOptions) *apperror.AppError {
	if opts.InterfaceName == "" {
		return apperror.NewValidationError("interface name is required")
	}

	if !opts.IsDHCP && opts.IP == "" {
		return apperror.NewValidationError("IP address is required for static configuration")
	}

	return nil
}

func (m *Manager) prepareSnapshot(
	ctx context.Context,
	ifaceName string,
) (Snapshot, *apperror.AppError) {
	currentRes := m.driver.GetInterface(ctx, ifaceName)
	if currentRes.IsFailure() {
		return Snapshot{}, currentRes.AppError()
	}

	snap := snapshotFromInterface(currentRes.Value)
	if persistErr := m.persistSnapshot(snap); persistErr != nil {
		return Snapshot{}, persistErr
	}

	return snap, nil
}

func buildDryRunResult(ifaceName, currentIP string) RollbackResultWrap {
	return result.Ok(RollbackResult{
		IsReverted:    false,
		InterfaceName: ifaceName,
		RestoredIP:    currentIP,
		Message:       "dry run completed successfully; no changes applied",
		IsValid:       true,
	})
}

func (m *Manager) executeChangeWithRollback(
	ctx context.Context,
	opts ChangeOptions,
	valOpts ValidationOptions,
	snap Snapshot,
) RollbackResultWrap {
	if applyErr := m.driver.ApplyConfig(ctx, opts); applyErr != nil {
		return m.orchestrateRollback(ctx, snap, "driver application failed: "+applyErr.Message)
	}

	if valErr := m.validateConnectivity(ctx, valOpts); valErr != nil {
		return m.orchestrateRollback(ctx, snap, "validation failed: "+valErr.Message)
	}

	return buildSuccessResult(opts.InterfaceName, opts.IP)
}

func (m *Manager) validateConnectivity(ctx context.Context, valOpts ValidationOptions) *apperror.AppError {
	if !valOpts.IsValidationActive {
		return nil
	}

	return m.validator.ValidateConnectivity(ctx, valOpts)
}

func buildSuccessResult(ifaceName, newIP string) RollbackResultWrap {
	return result.Ok(RollbackResult{
		IsReverted:    false,
		InterfaceName: ifaceName,
		RestoredIP:    "",
		Message:       "IP configuration applied successfully to " + newIP,
		IsValid:       true,
	})
}

func (m *Manager) orchestrateRollback(
	ctx context.Context,
	snap Snapshot,
	reason string,
) RollbackResultWrap {
	revertErr := m.driver.RevertConfig(ctx, snap)
	if revertErr != nil {
		msg := reason + " (rollback also failed: " + revertErr.Message + ")"
		return result.Fail[RollbackResult](apperror.NewExecutionError(msg))
	}

	return result.Ok(RollbackResult{
		IsReverted:    true,
		InterfaceName: snap.InterfaceName,
		RestoredIP:    snap.IP,
		Message:       reason + "; successfully rolled back to " + snap.IP,
		IsValid:       true,
	})
}

// RevertIP restores the network configuration from the latest snapshot.
func (m *Manager) RevertIP(ctx context.Context, ifaceName string) RollbackResultWrap {
	if m.driver == nil {
		return result.Fail[RollbackResult](apperror.NewExecutionError("network IP driver not available on this platform"))
	}

	if ifaceName == "" {
		return result.Fail[RollbackResult](apperror.NewValidationError("interface name is required"))
	}

	snap, err := m.loadLatestSnapshot(ifaceName)
	if err != nil {
		return result.Fail[RollbackResult](err)
	}

	return m.executeRevert(ctx, snap)
}

func (m *Manager) executeRevert(ctx context.Context, snap Snapshot) RollbackResultWrap {
	if err := m.driver.RevertConfig(ctx, snap); err != nil {
		return result.Fail[RollbackResult](err)
	}

	return result.Ok(RollbackResult{
		IsReverted:    true,
		InterfaceName: snap.InterfaceName,
		RestoredIP:    snap.IP,
		Message:       "configuration successfully reverted from snapshot to " + snap.IP,
		IsValid:       true,
	})
}

func (m *Manager) persistSnapshot(snap Snapshot) *apperror.AppError {
	rec := snapshotToRecord(snap)
	if m.db != nil {
		return m.db.InsertIPSnapshot(&rec)
	}

	return store.SaveIPSnapshotJSON(&rec)
}

func (m *Manager) loadLatestSnapshot(ifaceName string) (Snapshot, *apperror.AppError) {
	snap, isFound, err := m.loadSnapshotFromDB(ifaceName)
	if err != nil {
		return Snapshot{}, err
	}
	if isFound {
		return snap, nil
	}

	return loadFallbackSnapshot()
}

func (m *Manager) loadSnapshotFromDB(ifaceName string) (Snapshot, bool, *apperror.AppError) {
	if m.db == nil {
		return Snapshot{}, false, nil
	}

	rec, err := m.db.GetLatestIPSnapshot(ifaceName)
	if err != nil || rec == nil {
		return Snapshot{}, false, nil
	}

	return recordToSnapshot(*rec), true, nil
}

func loadFallbackSnapshot() (Snapshot, *apperror.AppError) {
	rec, err := store.LoadLatestIPSnapshotJSON()
	if err != nil {
		return Snapshot{}, err
	}

	return recordToSnapshot(*rec), nil
}

func snapshotFromInterface(info InterfaceInfo) Snapshot {
	return Snapshot{
		InterfaceName: info.Name,
		IP:            info.IP,
		Netmask:       info.Netmask,
		Gateway:       info.Gateway,
		DNS:           info.DNS,
		IsDHCP:        info.IsDHCP,
		Timestamp:     time.Now().Unix(),
		IsValid:       info.IsValid,
	}
}

func snapshotToRecord(snap Snapshot) store.IPSnapshotRecord {
	return store.IPSnapshotRecord{
		IPSnapshotId:  snap.IPSnapshotId,
		InterfaceName: snap.InterfaceName,
		IP:            snap.IP,
		Netmask:       snap.Netmask,
		Gateway:       snap.Gateway,
		DNS:           strings.Join(snap.DNS, ","),
		IsDHCP:        snap.IsDHCP,
		Timestamp:     snap.Timestamp,
		Notes:         snap.Notes,
		Comments:      snap.Comments,
	}
}

func recordToSnapshot(rec store.IPSnapshotRecord) Snapshot {
	return Snapshot{
		IPSnapshotId:  rec.IPSnapshotId,
		InterfaceName: rec.InterfaceName,
		IP:            rec.IP,
		Netmask:       rec.Netmask,
		Gateway:       rec.Gateway,
		DNS:           parseDNSList(rec.DNS),
		IsDHCP:        rec.IsDHCP,
		Timestamp:     rec.Timestamp,
		Notes:         rec.Notes,
		Comments:      rec.Comments,
		IsValid:       rec.IP != "",
	}
}

func parseDNSList(dnsStr string) []string {
	if dnsStr == "" {
		return nil
	}

	return strings.Split(dnsStr, ",")
}

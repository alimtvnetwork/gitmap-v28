package netip

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

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

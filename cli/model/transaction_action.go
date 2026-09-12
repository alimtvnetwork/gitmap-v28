// Package model — TransactionActionRecord (v23+).
//
// One row per typed reversible step inside a Transaction. The Forward
// and Reverse payloads are opaque JSON blobs whose schema is owned by
// the per-Kind handler (see gitmap/txn/action.go).
package model

// TransactionActionRecord is one TransactionAction row.
type TransactionActionRecord struct {
	ID            int64  `json:"id"`
	TransactionID int64  `json:"transactionId"`
	Seq           int64  `json:"seq"`
	Kind          string `json:"kind"`
	ForwardJSON   string `json:"forwardJson"`
	ReverseJSON   string `json:"reverseJson"`
	BackupRef     int64  `json:"backupRef,omitempty"`
	AppliedAt     int64  `json:"appliedAt"`
	RevertedAt    int64  `json:"revertedAt,omitempty"`
}

package enums

// BatchStatusType represents the outcome status of a batch clone item.
type BatchStatusType string

const (
	BatchStatusOK      BatchStatusType = "ok"
	BatchStatusSkipped BatchStatusType = "skipped"
	BatchStatusFailed  BatchStatusType = "failed"
)

// IsValid reports whether the batch status is recognized.
func (b BatchStatusType) IsValid() bool {
	switch b {
	case BatchStatusOK, BatchStatusSkipped, BatchStatusFailed:
		return true
	default:
		return false
	}
}

// String returns the string representation.
func (b BatchStatusType) String() string {
	return string(b)
}

// IsOK reports whether the item succeeded.
func (b BatchStatusType) IsOK() bool {
	return b == BatchStatusOK
}

// IsSkipped reports whether the item was skipped.
func (b BatchStatusType) IsSkipped() bool {
	return b == BatchStatusSkipped
}

// IsFailed reports whether the item failed.
func (b BatchStatusType) IsFailed() bool {
	return b == BatchStatusFailed
}

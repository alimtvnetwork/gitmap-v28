package fileutil

import "os"

// FileModeType represents standard octal file permissions as an enum.
type FileModeType uint32

const (
	FileModeTypeDefaultFile          FileModeType = 0644
	FileModeTypeDefaultDir           FileModeType = 0755
	FileModeTypeReadOnly             FileModeType = 0444
	FileModeTypeWriteOnly            FileModeType = 0222
	FileModeTypeReadWrite            FileModeType = 0666
	FileModeTypeOwnerOnly            FileModeType = 0700
	FileModeTypeOwnerReadWrite       FileModeType = 0600
	FileModeTypeOwnerGroupReadWrite  FileModeType = 0660
	FileModeTypeAllPermission        FileModeType = 0777
	FileModeTypeAllExecute           FileModeType = 0111
	FileModeTypeAllReadExecute       FileModeType = 0555
	FileModeTypeOwnerFullGroupReadEx FileModeType = 0755
)

// ToFileMode converts FileModeType to standard library os.FileMode.
func (m FileModeType) ToFileMode() os.FileMode {
	if m == 0 {
		return os.FileMode(FileModeTypeDefaultFile)
	}

	return os.FileMode(m)
}

// String implements fmt.Stringer.
func (m FileModeType) String() string {
	return m.ToFileMode().String()
}

// IsExecutable reports whether the permission includes any execution bit.
func (m FileModeType) IsExecutable() bool {
	return (m & 0111) != 0
}

// IsWritable reports whether the permission includes any write bit.
func (m FileModeType) IsWritable() bool {
	return (m & 0222) != 0
}

// IsReadOnly reports whether the permission has read but no write bits.
func (m FileModeType) IsReadOnly() bool {
	return (m&0444) != 0 && (m&0222) == 0
}

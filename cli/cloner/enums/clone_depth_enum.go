package enums

// CloneDepthType represents git clone depth and filter modes.
type CloneDepthType byte

const (
	CloneDepthFull CloneDepthType = iota
	CloneDepthShallow
	CloneDepthBlobless
	CloneDepthTreeless
)

var depthStrings = [...]string{
	CloneDepthFull:     "full",
	CloneDepthShallow:  "shallow",
	CloneDepthBlobless: "blobless",
	CloneDepthTreeless: "treeless",
}

// IsValid reports whether the depth mode is recognized.
func (d CloneDepthType) IsValid() bool {
	return d <= CloneDepthTreeless
}

// String returns the human-readable string representation.
func (d CloneDepthType) String() string {
	if int(d) < len(depthStrings) {
		return depthStrings[d]
	}

	return "unknown"
}

// IsShallow reports whether the depth is shallow.
func (d CloneDepthType) IsShallow() bool {
	return d == CloneDepthShallow
}

// IsBlobless reports whether the depth is blobless.
func (d CloneDepthType) IsBlobless() bool {
	return d == CloneDepthBlobless
}

// IsTreeless reports whether the depth is treeless.
func (d CloneDepthType) IsTreeless() bool {
	return d == CloneDepthTreeless
}

// IsFull reports whether the depth is full.
func (d CloneDepthType) IsFull() bool {
	return d == CloneDepthFull
}

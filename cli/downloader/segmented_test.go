package downloader

import (
	"testing"
)

func TestCalculateSegments_DefaultTenSplits(t *testing.T) {
	totalSize := int64(10 * 1024 * 1024)
	segs := CalculateSegments(totalSize, DefaultSegmentCount)
	if len(segs) != 10 {
		t.Fatalf("expected 10 segments, got %d", len(segs))
	}

	if segs[0].Start != 0 {
		t.Errorf("expected first segment start 0, got %d", segs[0].Start)
	}

	if segs[len(segs)-1].End != totalSize-1 {
		t.Errorf("expected last segment end %d, got %d", totalSize-1, segs[len(segs)-1].End)
	}
}

func TestCalculateSegments_SmallChunkBoundary(t *testing.T) {
	totalSize := int64(20 * 1024)
	segs := CalculateSegments(totalSize, DefaultSegmentCount)
	if len(segs) != 4 {
		t.Fatalf("expected 4 segments bounded by 5KB MinChunkSize, got %d", len(segs))
	}

	if segs[0].Start != 0 {
		t.Errorf("expected start 0, got %d", segs[0].Start)
	}

	if segs[len(segs)-1].End != totalSize-1 {
		t.Errorf("expected end %d, got %d", totalSize-1, segs[len(segs)-1].End)
	}
}

func TestCalculateSegments_ZeroOrNegativeSize(t *testing.T) {
	if segs := CalculateSegments(0, 10); segs != nil {
		t.Errorf("expected nil for 0 size, got %v", segs)
	}

	if segs := CalculateSegments(-100, 10); segs != nil {
		t.Errorf("expected nil for negative size, got %v", segs)
	}
}

func verifySegmentContiguity(t *testing.T, totalSize int64, segs []SegmentRange) {
	for i := 1; i < len(segs); i++ {
		if segs[i].Start != segs[i-1].End+1 {
			t.Errorf("gap between segment %d and %d: end=%d, start=%d", i-1, i, segs[i-1].End, segs[i].Start)
		}
	}
}

func TestCalculateSegments_Contiguity(t *testing.T) {
	sizes := []int64{100, 105, 50000, 1048576, 5242880}
	for _, size := range sizes {
		segs := CalculateSegments(size, 10)
		if len(segs) == 0 {
			t.Fatalf("expected non-empty segments for size %d", size)
		}

		verifySegmentContiguity(t, size, segs)
	}
}

func TestHasAria2c(t *testing.T) {
	hasAria := HasAria2c()
	t.Logf("aria2c detected in environment: %v", hasAria)
}

func TestSegmentRange_FieldEquivalence(t *testing.T) {
	seg := makeSegment(3, 100, 200)
	if seg.Start != seg.StartOffset || seg.End != seg.EndOffset {
		t.Errorf("offset fields do not match start/end: %+v", seg)
	}

	if seg.Index != 3 {
		t.Errorf("expected index 3, got %d", seg.Index)
	}
}

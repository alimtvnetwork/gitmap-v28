package openfiletype_test

import (
	"encoding/json"
	"os"
	"testing"

	"coding-guidelines/common/pkg/enum/openfiletype"
)

func TestInvalidZeroValue(t *testing.T) {
	var zero openfiletype.Variant
	if zero != openfiletype.Invalid {
		t.Fatalf("expected zero value to be Invalid, got %v", zero)
	}

	if zero.IsValid() {
		t.Fatalf("expected Invalid to not be valid")
	}

	if !zero.IsInvalid() {
		t.Fatalf("expected Invalid to be IsInvalid")
	}

	if zero.String() != "Invalid" {
		t.Fatalf("expected String 'Invalid', got %s", zero.String())
	}
}

func TestVariantsAndFlags(t *testing.T) {
	tests := []struct {
		variant       openfiletype.Variant
		expectedName  string
		expectedFlags int
	}{
		{openfiletype.ReadOnly, "ReadOnly", os.O_RDONLY},
		{openfiletype.WriteOnly, "WriteOnly", os.O_WRONLY},
		{openfiletype.ReadWrite, "ReadWrite", os.O_RDWR},
		{openfiletype.Append, "Append", os.O_WRONLY | os.O_APPEND},
		{openfiletype.CreateAppend, "CreateAppend", os.O_CREATE | os.O_WRONLY | os.O_APPEND},
		{openfiletype.CreateTruncate, "CreateTruncate", os.O_CREATE | os.O_WRONLY | os.O_TRUNC},
		{openfiletype.CreateNew, "CreateNew", os.O_CREATE | os.O_EXCL | os.O_WRONLY},
		{openfiletype.ReadOrCreateOnly, "ReadOrCreateOnly", os.O_RDONLY | os.O_CREATE},
		{openfiletype.WriteOrCreateOnly, "WriteOrCreateOnly", os.O_WRONLY | os.O_CREATE},
		{openfiletype.ReadWriteOrCreateOnly, "ReadWriteOrCreateOnly", os.O_RDWR | os.O_CREATE},
	}

	for _, tc := range tests {
		if tc.variant.Name() != tc.expectedName {
			t.Errorf("expected name %s, got %s", tc.expectedName, tc.variant.Name())
		}

		if tc.variant.Flags() != tc.expectedFlags {
			t.Errorf("expected flags %d, got %d for %s", tc.expectedFlags, tc.variant.Flags(), tc.expectedName)
		}

		if !tc.variant.IsValid() {
			t.Errorf("expected %s to be valid", tc.expectedName)
		}
	}
}

func TestCheckers(t *testing.T) {
	ro := openfiletype.ReadOnly
	if !ro.IsReadOnly() {
		t.Fatalf("expected IsReadOnly true")
	}

	if ro.IsWriteOnly() {
		t.Fatalf("expected IsWriteOnly false")
	}

	ca := openfiletype.CreateAppend
	if !ca.IsCreateAppend() {
		t.Fatalf("expected IsCreateAppend true")
	}

	ct := openfiletype.CreateTruncate
	if !ct.IsCreateTruncate() {
		t.Fatalf("expected IsCreateTruncate true")
	}

	cn := openfiletype.CreateNew
	if !cn.IsCreateNew() {
		t.Fatalf("expected IsCreateNew true")
	}

	roco := openfiletype.ReadOrCreateOnly
	if !roco.IsReadOrCreateOnly() {
		t.Fatalf("expected IsReadOrCreateOnly true")
	}

	woco := openfiletype.WriteOrCreateOnly
	if !woco.IsWriteOrCreateOnly() {
		t.Fatalf("expected IsWriteOrCreateOnly true")
	}

	rwoco := openfiletype.ReadWriteOrCreateOnly
	if !rwoco.IsReadWriteOrCreateOnly() {
		t.Fatalf("expected IsReadWriteOrCreateOnly true")
	}
}

func TestParse(t *testing.T) {
	res := openfiletype.Parse("createappend")
	if res.IsFailed() {
		t.Fatalf("expected parse to succeed: %v", res.Fault())
	}

	if res.Data() != openfiletype.CreateAppend {
		t.Fatalf("expected CreateAppend, got %v", res.Data())
	}

	resRO := openfiletype.Parse("readorcreateonly")
	if resRO.IsFailed() || resRO.Data() != openfiletype.ReadOrCreateOnly {
		t.Fatalf("expected ReadOrCreateOnly, got %v", resRO.Data())
	}

	resWO := openfiletype.Parse("WriteOrCreateOnly")
	if resWO.IsFailed() || resWO.Data() != openfiletype.WriteOrCreateOnly {
		t.Fatalf("expected WriteOrCreateOnly, got %v", resWO.Data())
	}

	resRWO := openfiletype.Parse("READWRITEORCREATEONLY")
	if resRWO.IsFailed() || resRWO.Data() != openfiletype.ReadWriteOrCreateOnly {
		t.Fatalf("expected ReadWriteOrCreateOnly, got %v", resRWO.Data())
	}

	badRes := openfiletype.Parse("invalid-mode-string")
	if badRes.IsSuccess() {
		t.Fatalf("expected failure on bad string")
	}
}

func TestAllAndValues(t *testing.T) {
	all := openfiletype.All()
	if len(all) != 10 {
		t.Fatalf("expected 10 valid variants, got %d", len(all))
	}

	values := openfiletype.Values()
	if len(values) != 10 {
		t.Fatalf("expected 10 string values, got %d", len(values))
	}
}

func TestJSONRoundtrip(t *testing.T) {
	val := openfiletype.CreateTruncate

	data, err := json.Marshal(val)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	if string(data) != `"CreateTruncate"` {
		t.Fatalf("unexpected json: %s", string(data))
	}

	var parsed openfiletype.Variant
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if parsed != openfiletype.CreateTruncate {
		t.Fatalf("expected CreateTruncate, got %v", parsed)
	}
}

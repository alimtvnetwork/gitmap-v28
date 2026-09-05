package dbengine

import (
	"testing"
)

func TestScanString(t *testing.T) {
	if ScanString(nil) != "" {
		t.Errorf("expected empty string for nil")
	}
	if ScanString("hello") != "hello" {
		t.Errorf("expected 'hello'")
	}
	if ScanString([]byte("world")) != "world" {
		t.Errorf("expected 'world'")
	}
	if ScanString(123) != "123" {
		t.Errorf("expected '123'")
	}
}

func TestScanInt(t *testing.T) {
	if ScanInt(nil) != 0 {
		t.Errorf("expected 0 for nil")
	}
	if ScanInt(42) != 42 {
		t.Errorf("expected 42 for int")
	}
	if ScanInt(int64(100)) != 100 {
		t.Errorf("expected 100 for int64")
	}
	if ScanInt(int32(50)) != 50 {
		t.Errorf("expected 50 for int32")
	}
	if ScanInt(uint64(25)) != 25 {
		t.Errorf("expected 25 for uint64")
	}
	if ScanInt("not-int") != 0 {
		t.Errorf("expected 0 for non-int")
	}
}

func TestScanInt64(t *testing.T) {
	if ScanInt64(nil) != 0 {
		t.Errorf("expected 0 for nil")
	}
	if ScanInt64(int64(4200)) != 4200 {
		t.Errorf("expected 4200 for int64")
	}
	if ScanInt64(42) != 42 {
		t.Errorf("expected 42 for int")
	}
	if ScanInt64(int32(50)) != 50 {
		t.Errorf("expected 50 for int32")
	}
	if ScanInt64(uint64(999)) != 999 {
		t.Errorf("expected 999 for uint64")
	}
	if ScanInt64("not-int") != 0 {
		t.Errorf("expected 0 for non-int")
	}
}

func TestScanUint64(t *testing.T) {
	if ScanUint64(nil) != 0 {
		t.Errorf("expected 0 for nil")
	}
	if ScanUint64(uint64(555)) != 555 {
		t.Errorf("expected 555 for uint64")
	}
	if ScanUint64(int64(666)) != 666 {
		t.Errorf("expected 666 for int64")
	}
	if ScanUint64(777) != 777 {
		t.Errorf("expected 777 for int")
	}
	if ScanUint64(uint(888)) != 888 {
		t.Errorf("expected 888 for uint")
	}
	if ScanUint64("invalid") != 0 {
		t.Errorf("expected 0 for invalid")
	}
}

func TestScanUint(t *testing.T) {
	if ScanUint(nil) != 0 {
		t.Errorf("expected 0 for nil")
	}
	if ScanUint(uint(123)) != 123 {
		t.Errorf("expected 123 for uint")
	}
	if ScanUint(int64(456)) != 456 {
		t.Errorf("expected 456 for int64")
	}
}

func TestScanBool(t *testing.T) {
	if ScanBool(nil) != false {
		t.Errorf("expected false for nil")
	}
	if ScanBool(true) != true {
		t.Errorf("expected true for true")
	}
	if ScanBool(false) != false {
		t.Errorf("expected false for false")
	}
	if ScanBool(int64(1)) != true {
		t.Errorf("expected true for int64(1)")
	}
	if ScanBool(int64(0)) != false {
		t.Errorf("expected false for int64(0)")
	}
	if ScanBool(1) != true {
		t.Errorf("expected true for int(1)")
	}
	if ScanBool(0) != false {
		t.Errorf("expected false for int(0)")
	}
	if ScanBool(uint64(1)) != true {
		t.Errorf("expected true for uint64(1)")
	}
	if ScanBool(uint64(0)) != false {
		t.Errorf("expected false for uint64(0)")
	}
	if ScanBool("not-bool") != false {
		t.Errorf("expected false for string")
	}
}

func TestScanFloat64(t *testing.T) {
	if ScanFloat64(nil) != 0.0 {
		t.Errorf("expected 0.0 for nil")
	}
	if ScanFloat64(3.14) != 3.14 {
		t.Errorf("expected 3.14 for float64")
	}
	if ScanFloat64(float32(2.5)) != 2.5 {
		t.Errorf("expected 2.5 for float32")
	}
	if ScanFloat64(int64(10)) != 10.0 {
		t.Errorf("expected 10.0 for int64")
	}
	if ScanFloat64(5) != 5.0 {
		t.Errorf("expected 5.0 for int")
	}
	if ScanFloat64("bad") != 0.0 {
		t.Errorf("expected 0.0 for bad")
	}
}

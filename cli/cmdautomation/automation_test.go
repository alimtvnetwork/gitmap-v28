package cmdautomation

import (
	"bytes"
	"testing"
	"time"
)

func TestRegexRegistry_LazyCompilation(t *testing.T) {
	pattern := `func\s+Test[A-Z]\w*`
	reg1, err1 := GetRegex(pattern, false)
	if err1 != nil {
		t.Fatalf("unexpected error compiling regex: %v", err1)
	}

	reg2, err2 := GetRegex(pattern, false)
	if err2 != nil {
		t.Fatalf("unexpected error retrieving regex: %v", err2)
	}

	if reg1 != reg2 {
		t.Error("expected identical pointer from lazy regex cache")
	}
}

func TestHasLiteralMatch(t *testing.T) {
	data := []byte("Hello Golang World\nNew Line Here")
	hasMatch := HasLiteralMatch(data, "golang", true)
	if !hasMatch {
		t.Error("expected case-insensitive match for 'golang'")
	}

	hasExactMatch := HasLiteralMatch(data, "golang", false)
	if hasExactMatch {
		t.Error("expected case-sensitive miss for 'golang'")
	}
}

func TestPolyglotExtensionDetection(t *testing.T) {
	exts := []string{".go", ".ts", ".tsx", ".rs", ".cs", ".java", ".py", ".md", ".json"}
	for _, ext := range exts {
		isText := IsPolyglotTextExtension(ext)
		if !isText {
			t.Errorf("expected %s to be recognized as polyglot text", ext)
		}
	}

	isBin := IsPolyglotTextExtension(".exe")
	if isBin {
		t.Error("expected .exe to not be polyglot text")
	}
}

func TestBinaryProbeGuard(t *testing.T) {
	textData := []byte("Hello clean text file without null bytes")
	isBinText := HasBinaryContent(textData)
	if isBinText {
		t.Error("expected text data to not be flagged as binary")
	}

	binData := []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x00, 0x01}
	isBin := HasBinaryContent(binData)
	if !isBin {
		t.Error("expected null byte data to be flagged as binary")
	}
}

func TestStripUtf8Bom(t *testing.T) {
	dataWithBom := []byte{0xEF, 0xBB, 0xBF, 'H', 'i'}
	stripped := StripUtf8Bom(dataWithBom)
	if string(stripped) != "Hi" {
		t.Errorf("expected 'Hi', got %q", string(stripped))
	}

	dataNoBom := []byte("NoBom")
	unmodified := StripUtf8Bom(dataNoBom)
	if string(unmodified) != "NoBom" {
		t.Errorf("expected 'NoBom', got %q", string(unmodified))
	}
}

func TestNormalizeContent_CrlfAndWhitespace(t *testing.T) {
	raw := []byte("line 1   \r\nline 2\t\r\nline 3\r\n\r\n\r\n")
	cleaned, crlfs, isModified := NormalizeContent(raw)

	if !isModified {
		t.Error("expected content to be flagged as modified")
	}
	if crlfs != 3 {
		t.Errorf("expected 3 CRLFs, got %d", crlfs)
	}

	expected := []byte("line 1\nline 2\nline 3\n")
	if !bytes.Equal(cleaned, expected) {
		t.Errorf("expected %q, got %q", string(expected), string(cleaned))
	}
}

func TestMemoryCache_SetGetClear(t *testing.T) {
	cache := GlobalCache()
	cache.Clear()

	cache.SetFile("test.go", []byte("package main"))
	data, isHit := cache.GetFile("test.go")
	if !isHit {
		t.Error("expected cache hit for test.go")
	}
	if string(data) != "package main" {
		t.Errorf("expected 'package main', got %q", string(data))
	}

	cache.Clear()
	_, isHitAfter := cache.GetFile("test.go")
	if isHitAfter {
		t.Error("expected cache miss after Clear()")
	}
}

func TestAssembleMetric(t *testing.T) {
	goDur := 10 * time.Millisecond
	pyDur := 100 * time.Millisecond
	metric := assembleMetric("Test Operation", goDur, pyDur)

	if metric.GoSpeedup < 9.9 || metric.GoSpeedup > 10.1 {
		t.Errorf("expected ~10x speedup, got %.2f", metric.GoSpeedup)
	}
	if metric.GoTempBytes != 0 {
		t.Errorf("expected 0 bytes for Go temp disk usage, got %d", metric.GoTempBytes)
	}
}

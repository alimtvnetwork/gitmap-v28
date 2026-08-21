package archive

import (
	"compress/flate"
	"compress/gzip"
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestFormatFromPath_DoubleExtensions(t *testing.T) {
	cases := map[string]Format{
		"archive.tar.gz":  FormatTarGz,
		"archive.tgz":     FormatTarGz,
		"archive.tar.bz2": FormatTarBz2,
		"archive.tbz2":    FormatTarBz2,
		"archive.tar.xz":  FormatTarXz,
		"archive.txz":     FormatTarXz,
		"archive.tar.zst": FormatTarZst,
		"archive.tzst":    FormatTarZst,
	}
	for path, want := range cases {
		if got := FormatFromPath(path); got != want {
			t.Errorf("FormatFromPath(%q) = %q, want %q", path, got, want)
		}
	}
}

func TestFormatFromPath_SingleExtensions(t *testing.T) {
	cases := map[string]Format{
		"archive.zip":     FormatZip,
		"archive.tar":     FormatTar,
		"archive.gz":      FormatGz,
		"archive.bz2":     FormatBz2,
		"archive.xz":      FormatXz,
		"archive.zst":     FormatZst,
		"archive.7z":      Format7z,
		"archive.rar":     FormatRar,
		"archive.unknown": FormatUnknown,
	}
	for path, want := range cases {
		if got := FormatFromPath(path); got != want {
			t.Errorf("FormatFromPath(%q) = %q, want %q", path, got, want)
		}
	}
}

func TestFormatExtension_RoundTrip(t *testing.T) {
	formats := []Format{
		FormatZip, FormatTar, FormatTarGz, FormatTarBz2,
		FormatTarXz, FormatTarZst, FormatGz, FormatBz2,
		FormatXz, FormatZst, Format7z, FormatRar,
	}
	for _, f := range formats {
		ext := f.Extension()
		if len(ext) == 0 {
			t.Errorf("Extension() for %q returned empty string", f)
		}
		if got := FormatFromPath("sample" + ext); got != f {
			t.Errorf("FormatFromPath(sample%s) = %q, want %q", ext, got, f)
		}
	}
}

func TestBuildArchiver_SupportedFormats(t *testing.T) {
	for _, f := range []Format{FormatZip, FormatTar, FormatTarGz, FormatTarBz2, FormatTarXz, FormatTarZst} {
		archiver, err := buildArchiver(f, ModeStandard)
		if err != nil {
			t.Errorf("buildArchiver(%q) returned unexpected error: %v", f, err)
		}
		if archiver == nil {
			t.Errorf("buildArchiver(%q) returned nil archiver", f)
		}
	}
}

func TestBuildArchiver_UnsupportedFormats(t *testing.T) {
	for _, f := range []Format{Format7z, FormatRar, FormatGz, FormatBz2, FormatXz, FormatZst, FormatUnknown} {
		_, err := buildArchiver(f, ModeStandard)
		if err == nil {
			t.Errorf("buildArchiver(%q) expected error, got nil", f)
		}
	}
}

func TestCompressionLevels(t *testing.T) {
	if gzipLevel(ModeFast) != gzip.BestSpeed || gzipLevel(ModeBest) != gzip.BestCompression {
		t.Errorf("gzipLevel mappings unexpected")
	}
	if bz2Level(ModeFast) != 1 || bz2Level(ModeBest) != 9 {
		t.Errorf("bz2Level mappings unexpected")
	}
	if flateLevel(ModeFast) != flate.BestSpeed || flateLevel(ModeBest) != flate.BestCompression {
		t.Errorf("flateLevel mappings unexpected")
	}
	if FlateLevelForMode(ModeStandard) != flate.DefaultCompression {
		t.Errorf("FlateLevelForMode unexpected")
	}
}

func TestMatchAny_Patterns(t *testing.T) {
	if matchAny("test.txt", nil, true) == false {
		t.Errorf("empty patterns should return emptyDefault=true")
	}
	if matchAny("test.txt", nil, false) {
		t.Errorf("empty patterns should return emptyDefault=false")
	}
	if matchAny("dir/test.txt", []string{"*.txt"}, false) == false {
		t.Errorf("basename match failed")
	}
	if matchAny("dir/test.txt", []string{"dir/*"}, false) == false {
		t.Errorf("path match failed")
	}
}

func TestCreateArchive_ZipBasic(t *testing.T) {
	tempDir := t.TempDir()
	srcFile := filepath.Join(tempDir, "hello.txt")
	if err := os.WriteFile(srcFile, []byte("hello world"), 0o644); err != nil {
		t.Fatal(err)
	}
	outPath := filepath.Join(tempDir, "out.zip")
	opts := CreateOptions{OutputPath: outPath, Sources: []string{srcFile}, Mode: ModeStandard}
	res, err := CreateArchive(context.Background(), opts)
	if err != nil {
		t.Fatalf("CreateArchive failed: %v", err)
	}
	if res.EntriesWritten != 1 || res.Format != FormatZip {
		t.Fatalf("unexpected CreateResult: %+v", res)
	}
}

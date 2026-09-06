package fileutil

import (
	"path/filepath"
	"testing"
)

type dummyData struct {
	Name string `json:"name" yaml:"name"`
	Age  int    `json:"age" yaml:"age"`
}

func TestParsersAndExporters(t *testing.T) {
	tmp := t.TempDir()

	t.Run("Text and Lines", func(t *testing.T) {
		path := filepath.Join(tmp, "lines.txt")
		lines := []string{"hello", "world"}
		
		expRes := ExportLines(path, lines, FilePermStandard)
		if expRes.HasError() {
			t.Fatalf("Failed to export lines: %v", expRes.Fault().Error())
		}

		readRes := ReadLines(path)
		if readRes.HasError() || len(readRes.Data()) != 2 {
			t.Fatalf("Failed to read lines correctly")
		}

		txtRes := ReadText(path)
		if txtRes.HasError() || txtRes.Data() != "hello\nworld\n" {
			t.Fatalf("Failed to read text correctly")
		}
	})

	t.Run("JSON", func(t *testing.T) {
		path := filepath.Join(tmp, "data.json")
		data := dummyData{Name: "Alice", Age: 30}

		expRes := ExportJSON(path, data, FilePermStandard)
		if expRes.HasError() {
			t.Fatalf("Failed to export JSON: %v", expRes.Fault().Error())
		}

		readRes := ReadJSON[dummyData](path)
		if readRes.HasError() || readRes.Data().Name != "Alice" {
			t.Fatalf("Failed to read JSON correctly")
		}
	})

	t.Run("YAML", func(t *testing.T) {
		path := filepath.Join(tmp, "data.yaml")
		data := dummyData{Name: "Bob", Age: 40}

		expRes := ExportYAML(path, data, FilePermStandard)
		if expRes.HasError() {
			t.Fatalf("Failed to export YAML: %v", expRes.Fault().Error())
		}

		readRes := ReadYAML[dummyData](path)
		if readRes.HasError() || readRes.Data().Name != "Bob" {
			t.Fatalf("Failed to read YAML correctly")
		}
	})
}

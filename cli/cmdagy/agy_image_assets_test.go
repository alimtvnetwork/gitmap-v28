package cmdagy

import (
	"bytes"
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const samplePngBase64 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII="

func TestExtractAndStoreBase64Images_SinglePNG(t *testing.T) {
	tempDir := t.TempDir()
	prompt := "Please inspect:\n![Screenshot](data:image/png;base64," + samplePngBase64 + ")\nDone."

	newPrompt, savedFiles, err := ExtractAndStoreBase64Images(prompt, tempDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertSingleExtractionResults(t, prompt, newPrompt, savedFiles, tempDir)
}

func assertSingleExtractionResults(t *testing.T, orig, newPrompt string, savedFiles []string, tempDir string) {
	if len(savedFiles) != 1 {
		t.Fatalf("expected 1 saved file, got %d", len(savedFiles))
	}

	verifySavedFileContent(t, savedFiles[0])
	verifyRewrittenPrompt(t, newPrompt)
}

func verifySavedFileContent(t *testing.T, filePath string) {
	raw, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed reading saved file: %v", err)
	}

	expected, _ := base64.StdEncoding.DecodeString(samplePngBase64)
	if !bytes.Equal(raw, expected) {
		t.Errorf("saved content does not match decoded base64")
	}
}

func verifyRewrittenPrompt(t *testing.T, prompt string) {
	if strings.Contains(prompt, "data:image/png;base64") {
		t.Errorf("prompt still contains base64 string")
	}

	if !strings.Contains(prompt, "http://127.0.0.1:18188/assets/img-") {
		t.Errorf("prompt missing standalone asset URL: %s", prompt)
	}

	if !strings.Contains(prompt, "<!-- asset: ") {
		t.Errorf("prompt missing asset comment: %s", prompt)
	}
}

func TestExtractAndStoreBase64Images_NoImages(t *testing.T) {
	tempDir := t.TempDir()
	orig := "Simple prompt without any embedded images."

	newPrompt, savedFiles, err := ExtractAndStoreBase64Images(orig, tempDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if newPrompt != orig || len(savedFiles) != 0 {
		t.Errorf("expected unmodified prompt and zero saved files")
	}
}

func TestExtractAndStoreBase64Images_MultipleImages(t *testing.T) {
	tempDir := t.TempDir()
	prompt := "First: data:image/png;base64," + samplePngBase64 + "\nSecond: data:image/jpeg;base64,/9j/4AAQSkZJRg=="

	newPrompt, savedFiles, err := ExtractAndStoreBase64Images(prompt, tempDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(savedFiles) != 2 {
		t.Fatalf("expected 2 saved files, got %d", len(savedFiles))
	}

	if strings.Contains(newPrompt, "data:image/") {
		t.Errorf("rewritten prompt still contains data URIs")
	}
}

func TestFormatImageMarkdownStandalone(t *testing.T) {
	imagePath := filepath.Join(".ai-memory", "assets", "img-abcdef1234567890.png")
	serverURL := "http://127.0.0.1:18188"

	result := FormatImageMarkdownStandalone(imagePath, serverURL)
	expectedURL := "http://127.0.0.1:18188/assets/img-abcdef1234567890.png"
	expectedComment := "<!-- asset: .ai-memory/assets/img-abcdef1234567890.png -->"

	if !strings.Contains(result, expectedURL) {
		t.Errorf("expected markdown to contain %q, got %q", expectedURL, result)
	}

	if !strings.Contains(result, expectedComment) {
		t.Errorf("expected markdown to contain %q, got %q", expectedComment, result)
	}
}

func TestFormatImageMarkdownStandalone_Empty(t *testing.T) {
	if res := FormatImageMarkdownStandalone("", ""); res != "" {
		t.Errorf("expected empty string for blank imagePath, got %q", res)
	}
}

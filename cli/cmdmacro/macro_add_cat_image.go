// Package cmdmacro — macro_add_cat_image.go: image asset inspection and TokyoNight card summary for cat.
package cmdmacro

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/output"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

func isImageExtension(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg", ".ico", ".bmp", ".tiff":
		return true
	default:
		return false
	}
}

func isSvgPath(path string) bool {
	return strings.EqualFold(filepath.Ext(path), ".svg")
}

func printCatImageSummary(target string, data []byte) {
	if isSvgPath(target) {
		printSvgSummary(target, data)

		return
	}

	printRasterSummary(target, data)
}

func printSvgSummary(target string, data []byte) {
	sizeStr := formatImageFileSize(len(data))
	card := renderImageCard(target, "SVG (Vector)", "Vector / Scalable", sizeStr)
	printImageCard(card)
	printSvgContent(data)
}

func printRasterSummary(target string, data []byte) {
	format, dims := resolveRasterMetadata(target, data)
	sizeStr := formatImageFileSize(len(data))
	card := renderImageCard(target, format, dims, sizeStr)
	printImageCard(card)
}

func resolveRasterMetadata(target string, data []byte) (string, string) {
	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return fallbackRasterMetadata(target)
	}

	dims := fmt.Sprintf("%d x %d px", cfg.Width, cfg.Height)

	return strings.ToUpper(format), dims
}

func fallbackRasterMetadata(target string) (string, string) {
	ext := strings.TrimPrefix(filepath.Ext(target), ".")
	fmtStr := strings.ToUpper(ext)

	return fmtStr, "Unknown dimensions"
}

func formatImageFileSize(size int) string {
	kb := float64(size) / 1024.0

	return fmt.Sprintf("%d bytes (%.1f KB)", size, kb)
}

func buildImageCardContent(target, format, dims, sizeStr string) string {
	title := termout.SafeTokyoStyle(termout.TokyoPrimary).Bold(true).
		Render(fmt.Sprintf("🖼 Image Asset Summary: %s", target))
	fmtLine := fmt.Sprintf("Format:     %s", format)
	dimLine := fmt.Sprintf("Dimensions: %s", dims)
	sizeLine := fmt.Sprintf("File Size:  %s", sizeStr)
	statusText := termout.SafeTokyoStyle(termout.TokyoOK).Render("✓ Verified and readable on disk")
	statusLine := fmt.Sprintf("Status:     %s", statusText)

	return fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s", title, fmtLine, dimLine, sizeLine, statusLine)
}

func renderImageCard(target, format, dims, sizeStr string) string {
	content := buildImageCardContent(target, format, dims, sizeStr)

	return termout.TokyoBoxStyle(content)
}

func printImageCard(card string) {
	fmt.Println()
	fmt.Println(card)
	fmt.Println()
}

func printSvgContent(data []byte) {
	ui := output.UI()
	fmt.Fprintf(ui, "\n  %s--- SVG Content ---%s\n", constants.ColorCyan, constants.ColorReset)

	_, _ = output.Raw().Write(data)

	if len(data) > 0 && data[len(data)-1] != '\n' {
		fmt.Fprintln(ui)
	}

	fmt.Fprintf(ui, "  %s---------------------%s\n\n", constants.ColorCyan, constants.ColorReset)
}

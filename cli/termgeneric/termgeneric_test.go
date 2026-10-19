package termgeneric

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestColor_Wrappers(t *testing.T) {
	SetColorEnabled(true)
	defer SetColorEnabled(true)

	tests := []struct {
		name     string
		fn       func(string) string
		prefix   string
		expected string
	}{
		{"Red", Red, CodeRed, CodeRed + "sample" + CodeReset},
		{"Green", Green, CodeGreen, CodeGreen + "sample" + CodeReset},
		{"Cyan", Cyan, CodeCyan, CodeCyan + "sample" + CodeReset},
		{"Yellow", Yellow, CodeYellow, CodeYellow + "sample" + CodeReset},
		{"Magenta", Magenta, CodeMagenta, CodeMagenta + "sample" + CodeReset},
		{"Blue", Blue, CodeBlue, CodeBlue + "sample" + CodeReset},
		{"Dim", Dim, CodeDim, CodeDim + "sample" + CodeReset},
		{"Bold", Bold, CodeBold, CodeBold + "sample" + CodeReset},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.fn("sample")
			if res != tt.expected {
				t.Errorf("%s(\"sample\") = %q, want %q", tt.name, res, tt.expected)
			}
		})
	}
}

func TestColor_Toggle(t *testing.T) {
	SetColorEnabled(false)
	defer SetColorEnabled(true)

	if IsColorEnabled() {
		t.Errorf("expected IsColorEnabled() to be false")
	}

	red := Red("plain text")
	if red != "plain text" {
		t.Errorf("expected plain text when color disabled, got %q", red)
	}

	reset := Reset("plain text")
	if reset != "plain text" {
		t.Errorf("expected plain text for Reset when disabled, got %q", reset)
	}
}

func TestColor_StripAnsi(t *testing.T) {
	SetColorEnabled(true)

	sample := Red("colored") + " and " + Bold(Cyan("bold cyan"))
	stripped := StripAnsi(sample)
	if stripped != "colored and bold cyan" {
		t.Errorf("StripAnsi failed, got %q", stripped)
	}

	csi := "\r\033[2Kloading"
	strippedCsi := StripAnsi(csi)
	if strippedCsi != "\rloading" {
		t.Errorf("StripAnsi failed on CSI, got %q", strippedCsi)
	}

	plain := "pure text"
	if StripAnsi(plain) != plain {
		t.Errorf("StripAnsi altered plain text")
	}

	if StripAnsi("") != "" {
		t.Errorf("StripAnsi empty string failed")
	}
}

func TestHeader_RenderDefault(t *testing.T) {
	box := RenderBox("GitMap Terminal", 40)
	lines := strings.Split(box, "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}

	if !strings.HasPrefix(lines[0], "╭") || !strings.HasSuffix(lines[0], "╮") {
		t.Errorf("invalid top border: %s", lines[0])
	}

	if !strings.Contains(lines[1], "GitMap Terminal") {
		t.Errorf("missing title in middle line: %s", lines[1])
	}

	if !strings.HasPrefix(lines[2], "╰") || !strings.HasSuffix(lines[2], "╯") {
		t.Errorf("invalid bottom border: %s", lines[2])
	}

	w0 := VisibleWidth(lines[0])
	w1 := VisibleWidth(lines[1])
	w2 := VisibleWidth(lines[2])
	if w0 != 40 || w1 != 40 || w2 != 40 {
		t.Errorf("width mismatch: %d, %d, %d; want 40", w0, w1, w2)
	}
}

func TestHeader_RenderWithSubtitleAndAlign(t *testing.T) {
	cfg := HeaderConfig{
		Title:     "Repository Status",
		Subtitle:  "Branch: main | Clean",
		Width:     50,
		Border:    BorderSharp,
		Alignment: AlignCenter,
		Padding:   2,
	}

	rendered := RenderHeader(cfg)
	lines := strings.Split(rendered, "\n")
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines with subtitle, got %d", len(lines))
	}

	if !strings.HasPrefix(lines[0], "┌") || !strings.HasSuffix(lines[0], "┐") {
		t.Errorf("sharp border top failed: %s", lines[0])
	}
	if !strings.HasPrefix(lines[3], "└") || !strings.HasSuffix(lines[3], "┘") {
		t.Errorf("sharp border bottom failed: %s", lines[3])
	}

	for i, l := range lines {
		w := VisibleWidth(l)
		if w != 50 {
			t.Errorf("line %d width %d != 50", i, w)
		}
	}
}

func TestHeader_Truncation(t *testing.T) {
	longTitle := "This is an extremely long title that definitely exceeds the allotted box width"
	truncated := TruncateTitle(longTitle, 20)
	visWidth := VisibleWidth(truncated)
	if visWidth > 20 {
		t.Errorf("TruncateTitle visible width %d exceeds 20", visWidth)
	}
	if !strings.HasSuffix(truncated, "...") {
		t.Errorf("expected truncated string to end with ..., got %s", truncated)
	}

	box := RenderBox(longTitle, 30)
	lines := strings.Split(box, "\n")
	for i, l := range lines {
		w := VisibleWidth(l)
		if w != 30 {
			t.Errorf("line %d width %d != 30", i, w)
		}
	}
}

func TestHeader_Styles(t *testing.T) {
	boxDouble := RenderHeader(HeaderConfig{
		Title:  "Double Border",
		Width:  35,
		Border: BorderDouble,
	})
	if !strings.HasPrefix(boxDouble, "╔") {
		t.Errorf("expected double top border, got %s", boxDouble)
	}

	boxAscii := RenderHeader(HeaderConfig{
		Title:  "ASCII Border",
		Width:  35,
		Border: BorderAscii,
	})
	if !strings.HasPrefix(boxAscii, "+") {
		t.Errorf("expected ascii top border, got %s", boxAscii)
	}

	subBox := RenderBoxWithSubtitle("Main Title", "Sub Title", 40)
	if !strings.Contains(subBox, "Main Title") || !strings.Contains(subBox, "Sub Title") {
		t.Errorf("RenderBoxWithSubtitle missing title or subtitle")
	}
}

func TestProgress_Spinner(t *testing.T) {
	var buf bytes.Buffer
	spin := NewSpinner("Scanning repos", &buf)

	frame1 := spin.Tick()
	if !strings.Contains(frame1, "Scanning repos") {
		t.Errorf("spinner frame missing label: %s", frame1)
	}

	spin.SetLabel("Updating index")
	frame2 := spin.Tick()
	if !strings.Contains(frame2, "Updating index") {
		t.Errorf("spinner updated label failed: %s", frame2)
	}

	spin.Clear()
	if !strings.Contains(buf.String(), LineClearSequence) {
		t.Errorf("Clear did not output LineClearSequence")
	}

	buf.Reset()
	spin.Success("All repositories processed")
	if !strings.Contains(buf.String(), "✔ All repositories processed") {
		t.Errorf("Success line missing checkmark or text: %s", buf.String())
	}

	buf.Reset()
	spin.Fail("Failed to sync")
	if !strings.Contains(buf.String(), "✖ Failed to sync") {
		t.Errorf("Fail line missing cross or text: %s", buf.String())
	}
}

func TestProgress_ProgressBar(t *testing.T) {
	var buf bytes.Buffer
	pb := NewProgressBar(100, "Downloading", &buf)

	pb.SetCurrent(25)
	if pb.Current() != 25 {
		t.Errorf("Current() = %d, want 25", pb.Current())
	}

	pct := pb.Percent()
	if pct != 25.0 {
		t.Errorf("Percent() = %f, want 25.0", pct)
	}

	pb.Increment()
	if pb.Current() != 26 {
		t.Errorf("Increment failed, current = %d", pb.Current())
	}

	pb.Add(24)
	if pb.Current() != 50 {
		t.Errorf("Add failed, current = %d", pb.Current())
	}

	rendered := pb.Render()
	if !strings.Contains(rendered, "50%") || !strings.Contains(rendered, "(50/100)") {
		t.Errorf("Render() unexpected content: %s", rendered)
	}

	pb.Draw()
	if !strings.Contains(buf.String(), LineClearSequence) {
		t.Errorf("Draw did not write LineClearSequence")
	}

	buf.Reset()
	pb.Finish("Complete!")
	if !strings.Contains(buf.String(), "Complete!") {
		t.Errorf("Finish did not print final message")
	}
}

func TestWorkerPool_ConcurrentExecution(t *testing.T) {
	ctx := context.Background()
	handler := func(ctx context.Context, item int) (int, error) {
		time.Sleep(5 * time.Millisecond)
		return item * 2, nil
	}

	pool := NewWorkerPool(ctx, 4, 10, handler)
	pool.Start()

	numTasks := 20
	for i := 1; i <= numTasks; i++ {
		isSubmitted := pool.Submit(i)
		if !isSubmitted {
			t.Fatalf("failed to submit task %d", i)
		}
	}

	go pool.Close()

	results := make(map[int]int)
	for res := range pool.Results() {
		if res.Err != nil {
			t.Errorf("unexpected task error: %v", res.Err)
		}
		results[res.Input] = res.Output
	}

	if len(results) != numTasks {
		t.Errorf("expected %d results, got %d", numTasks, len(results))
	}

	for i := 1; i <= numTasks; i++ {
		if results[i] != i*2 {
			t.Errorf("task %d result %d != %d", i, results[i], i*2)
		}
	}

	if pool.HasErrors() {
		t.Errorf("expected no errors in pool")
	}
}

func TestWorkerPool_ErrorsAndCancellation(t *testing.T) {
	ctx := context.Background()
	errSample := errors.New("simulated failure")
	handler := func(ctx context.Context, item string) (string, error) {
		if item == "fail" {
			return "", errSample
		}
		return "ok:" + item, nil
	}

	items := []string{"apple", "fail", "banana"}
	results, err := RunBatch(ctx, items, 2, handler)

	if len(results) != 3 {
		t.Errorf("expected 3 batch results, got %d", len(results))
	}
	if err == nil {
		t.Errorf("expected batch execution to return error")
	}

	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()

	pool := NewWorkerPool(cancelCtx, 2, 5, handler)
	pool.Start()

	isQueued := pool.Submit("test")
	if isQueued {
		t.Errorf("expected Submit to return false on canceled context")
	}

	pool.Stop()
}

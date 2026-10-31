package cmdcdp

import (
	"bytes"
	"encoding/json"
	"net/url"
	"testing"
)

func TestMaskPayloadRoundtrip(t *testing.T) {
	payload := []byte("hello websocket cdp world")
	mask := [4]byte{0x12, 0x34, 0x56, 0x78}
	masked := maskPayload(payload, mask)
	unmasked := maskPayload(masked, mask)

	if !bytes.Equal(payload, unmasked) {
		t.Fatalf("expected payload %s, got %s", payload, unmasked)
	}
}

func TestBuildFrameHeaderSmall(t *testing.T) {
	hdr := buildFrameHeader(50)
	if len(hdr) != 2 || hdr[0] != 0x81 || hdr[1] != (0x80|50) {
		t.Fatalf("unexpected header for length 50: %v", hdr)
	}
}

func TestBuildFrameHeaderMedium(t *testing.T) {
	hdr := buildFrameHeader(1000)
	if len(hdr) != 4 || hdr[0] != 0x81 || hdr[1] != (0x80|126) {
		t.Fatalf("unexpected header for length 1000: %v", hdr)
	}
}

func TestBuildFrameHeaderLarge(t *testing.T) {
	hdr := buildFrameHeader(70000)
	if len(hdr) != 10 || hdr[0] != 0x81 || hdr[1] != (0x80|127) {
		t.Fatalf("unexpected header for length 70000: %v", hdr)
	}
}

func TestResolveHostPortCustom(t *testing.T) {
	u, _ := url.Parse("ws://127.0.0.1:9222/devtools/page/123")
	if resolveHostPort(u) != "127.0.0.1:9222" {
		t.Fatalf("unexpected hostport: %s", resolveHostPort(u))
	}
}

func TestResolveHostPortDefault(t *testing.T) {
	uWS, _ := url.Parse("ws://localhost/devtools/browser")
	if resolveHostPort(uWS) != "localhost:80" {
		t.Fatalf("unexpected hostport: %s", resolveHostPort(uWS))
	}

	uWSS, _ := url.Parse("wss://remote-host/devtools/browser")
	if resolveHostPort(uWSS) != "remote-host:443" {
		t.Fatalf("unexpected hostport: %s", resolveHostPort(uWSS))
	}
}

func TestLocatePageTarget(t *testing.T) {
	targets := []CDPTarget{
		{ID: "1", Type: "background_page", WebSocketDebuggerURL: "ws://1"},
		{ID: "2", Type: "page", WebSocketDebuggerURL: "ws://2", Title: "Main Page"},
	}

	target := locatePageTarget(targets)
	if target == nil || target.ID != "2" {
		t.Fatalf("expected target 2, got %v", target)
	}
}

func TestUnpackEvalResult(t *testing.T) {
	raw := json.RawMessage(`{"result":{"type":"string","value":"test-value"}}`)
	unpacked := unpackEvalResult(raw)

	if string(unpacked) != `"test-value"` {
		t.Fatalf("unexpected eval unpacked: %s", string(unpacked))
	}
}

func TestExtractSubcommand(t *testing.T) {
	sub, rest := extractSubcommand([]string{"probe", "--port", "9222"})
	if sub != "probe" || len(rest) != 2 {
		t.Fatalf("unexpected subcmd: %s, rest: %v", sub, rest)
	}

	flagSub, flagRest := extractSubcommand([]string{"--port", "9222", "probe"})
	if flagSub != "" || len(flagRest) != 3 {
		t.Fatalf("unexpected flag subcmd: %s", flagSub)
	}
}

package transport

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func TestClassify(t *testing.T) {
	cases := map[string]string{
		"https://github.com/x/y.git":   constants.ScanTransportHTTPS,
		"ssh://git@github.com/x/y.git": constants.ScanTransportSSH,
		"git@github.com:x/y.git":       constants.ScanTransportSSH,
		"":                             constants.ScanTransportOther,
		"http://x/y.git":               constants.ScanTransportOther,
		"git://x/y.git":                constants.ScanTransportOther,
	}
	for in, want := range cases {
		if got := Classify(in); got != want {
			t.Errorf("Classify(%q) = %q, want %q", in, got, want)
		}
	}
}

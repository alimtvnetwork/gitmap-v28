package osclean

import (
	"testing"
)

func TestIsAntigravityProtected_Settings(t *testing.T) {
	protectedPaths := []string{
		`C:\Users\Admin\.gemini\config\config.json`,
		`C:\Users\Admin\.gemini\config\projects\040097de.json`,
		`C:\Users\Admin\.gemini\config\prompts\review.md`,
		`C:\Users\Admin\.gemini\config\plugins\my-plugin\manifest.json`,
		`C:\Users\Admin\.gemini\antigravity\conversations\conv123\transcript.jsonl`,
		`C:\Users\Admin\.gemini\antigravity\conversation_summaries.db`,
		`C:\Users\Admin\.gemini\antigravity\conversation_summaries.db-wal`,
		`C:\Users\Admin\.gemini\antigravity\agyhub_summaries_proto.pb`,
		`C:\Users\Admin\.gemini\antigravity\installation_id`,
		`C:\Users\Admin\.gemini\antigravity\antigravity_state.pbtxt`,
		`C:\Users\Admin\.gemini\antigravity\knowledge\doc.md`,
		`C:\Users\Admin\.gemini\antigravity\builtin\skills\guide.md`,
		`C:\Users\Admin\AppData\Roaming\Antigravity\Preferences`,
		`C:\Users\Admin\AppData\Roaming\Antigravity\Local State`,
		`C:\Users\Admin\AppData\Roaming\Antigravity\app_storage.json`,
		`C:\Users\Admin\AppData\Roaming\Antigravity\Local Storage\leveldb\000001.ldb`,
		`C:\Users\Admin\AppData\Roaming\Antigravity\Session Storage\000001.log`,
		`C:\Users\Admin\AppData\Roaming\Antigravity\Network\Cookies`,
		`C:\Users\Admin\AppData\Roaming\Antigravity\Network\Network Persistent State`,
		`C:\Users\Admin\AppData\Roaming\Antigravity\SharedStorage`,
		`C:\Users\Admin\AppData\Roaming\Antigravity\bin\agy-node.cmd`,
		`C:\Users\Admin\.gemini\config`,
		`C:\Users\Admin\AppData\Roaming\Antigravity`,
	}

	for _, p := range protectedPaths {
		if !IsAntigravityProtected(p) {
			t.Errorf("expected path to be protected: %s", p)
		}
	}
}

func TestIsAntigravityProtected_SafeCaches(t *testing.T) {
	safeCachePaths := []string{
		`C:\Users\Admin\AppData\Roaming\Antigravity\Cache\data_0`,
		`C:\Users\Admin\AppData\Roaming\Antigravity\Code Cache\js\script.bin`,
		`C:\Users\Admin\AppData\Roaming\Antigravity\GPUCache\data_1`,
		`C:\Users\Admin\AppData\Roaming\Antigravity\DawnWebGPUCache\index`,
		`C:\Users\Admin\AppData\Roaming\Antigravity\blob_storage\entry.bin`,
		`C:\Users\Admin\AppData\Roaming\Antigravity\logs\language_server.log`,
		`C:\Users\Admin\.gemini\antigravity\crashes\crash_123.log`,
		`C:\Users\Admin\.gemini\antigravity\scratch\scratch_001.tmp`,
		`C:\Users\Admin\AppData\Local\antigravity-updater\installer.exe`,
		`C:\dev-tool\go\pkg\mod\module.zip`,
		`C:\Users\Admin\AppData\Local\npm-cache\tarball.tgz`,
	}

	for _, p := range safeCachePaths {
		if IsAntigravityProtected(p) {
			t.Errorf("expected safe cache path NOT to be protected: %s", p)
		}
	}
}

func TestIsAntigravityProtected_ProtectedFileInScratch(t *testing.T) {
	scratchProtected := `C:\Users\Admin\.gemini\antigravity\scratch\config.json`
	if !IsAntigravityProtected(scratchProtected) {
		t.Errorf("expected config.json inside scratch to remain protected")
	}
}

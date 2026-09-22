package osclean

import (
	"testing"
)

func TestIsAntigravityProtected_Settings(t *testing.T) {
	protectedPaths := []string{
		`D:\mockhome\.gemini\config\config.json`,
		`D:\mockhome\.gemini\config\projects\040097de.json`,
		`D:\mockhome\.gemini\config\prompts\review.md`,
		`D:\mockhome\.gemini\config\plugins\my-plugin\manifest.json`,
		`D:\mockhome\.gemini\antigravity\conversations\conv123\transcript.jsonl`,
		`D:\mockhome\.gemini\antigravity\conversation_summaries.db`,
		`D:\mockhome\.gemini\antigravity\conversation_summaries.db-wal`,
		`D:\mockhome\.gemini\antigravity\agyhub_summaries_proto.pb`,
		`D:\mockhome\.gemini\antigravity\installation_id`,
		`D:\mockhome\.gemini\antigravity\antigravity_state.pbtxt`,
		`D:\mockhome\.gemini\antigravity\knowledge\doc.md`,
		`D:\mockhome\.gemini\antigravity\builtin\skills\guide.md`,
		`D:\mockhome\AppData\Roaming\Antigravity\Preferences`,
		`D:\mockhome\AppData\Roaming\Antigravity\Local State`,
		`D:\mockhome\AppData\Roaming\Antigravity\app_storage.json`,
		`D:\mockhome\AppData\Roaming\Antigravity\Local Storage\leveldb\000001.ldb`,
		`D:\mockhome\AppData\Roaming\Antigravity\Session Storage\000001.log`,
		`D:\mockhome\AppData\Roaming\Antigravity\Network\Cookies`,
		`D:\mockhome\AppData\Roaming\Antigravity\Network\Network Persistent State`,
		`D:\mockhome\AppData\Roaming\Antigravity\SharedStorage`,
		`D:\mockhome\AppData\Roaming\Antigravity\bin\agy-node.cmd`,
		`D:\mockhome\.gemini\config`,
		`D:\mockhome\AppData\Roaming\Antigravity`,
	}

	for _, p := range protectedPaths {
		if !IsAntigravityProtected(p) {
			t.Errorf("expected path to be protected: %s", p)
		}
	}
}

func TestIsAntigravityProtected_SafeCaches(t *testing.T) {
	safeCachePaths := []string{
		`D:\mockhome\AppData\Roaming\Antigravity\Cache\data_0`,
		`D:\mockhome\AppData\Roaming\Antigravity\Code Cache\js\script.bin`,
		`D:\mockhome\AppData\Roaming\Antigravity\GPUCache\data_1`,
		`D:\mockhome\AppData\Roaming\Antigravity\DawnWebGPUCache\index`,
		`D:\mockhome\AppData\Roaming\Antigravity\blob_storage\entry.bin`,
		`D:\mockhome\AppData\Roaming\Antigravity\logs\language_server.log`,
		`D:\mockhome\.gemini\antigravity\crashes\crash_123.log`,
		`D:\mockhome\.gemini\antigravity\scratch\scratch_001.tmp`,
		`D:\mockhome\AppData\Local\antigravity-updater\installer.exe`,
		`C:\dev-tool\go\pkg\mod\module.zip`,
		`D:\mockhome\AppData\Local\npm-cache\tarball.tgz`,
	}

	for _, p := range safeCachePaths {
		if IsAntigravityProtected(p) {
			t.Errorf("expected safe cache path NOT to be protected: %s", p)
		}
	}
}

func TestIsAntigravityProtected_ProtectedFileInScratch(t *testing.T) {
	scratchProtected := `D:\mockhome\.gemini\antigravity\scratch\config.json`
	if !IsAntigravityProtected(scratchProtected) {
		t.Errorf("expected config.json inside scratch to remain protected")
	}
}

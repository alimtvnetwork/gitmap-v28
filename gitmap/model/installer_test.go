package model

import (
	"encoding/json"
	"testing"
)

func TestInstallerModel(t *testing.T) {
	t.Parallel()

	scriptRecord := InstallerScript{
		ID:           1,
		Name:         "NodeJS Installer",
		Slug:         "nodejs-installer",
		Description:  "Installs NodeJS LTS",
		TargetOS:     "win",
		Version:      "1.0.0",
		Instructions: `{"steps":["download","install"]}`,
		CreatedAt:    "2026-08-26T00:00:00Z",
		UpdatedAt:    "2026-08-26T00:00:00Z",
	}

	marshaledScript, scriptErr := json.Marshal(scriptRecord)
	if scriptErr != nil {
		t.Fatalf("json.Marshal(InstallerScript) failed: %v", scriptErr)
	}

	var unmarshaledScript InstallerScript
	unmarshalScriptErr := json.Unmarshal(marshaledScript, &unmarshaledScript)
	if unmarshalScriptErr != nil {
		t.Fatalf("json.Unmarshal(InstallerScript) failed: %v", unmarshalScriptErr)
	}

	isScriptEqual := unmarshaledScript.ID == scriptRecord.ID &&
		unmarshaledScript.Name == scriptRecord.Name &&
		unmarshaledScript.Slug == scriptRecord.Slug &&
		unmarshaledScript.Description == scriptRecord.Description &&
		unmarshaledScript.TargetOS == scriptRecord.TargetOS &&
		unmarshaledScript.Version == scriptRecord.Version &&
		unmarshaledScript.Instructions == scriptRecord.Instructions &&
		unmarshaledScript.CreatedAt == scriptRecord.CreatedAt &&
		unmarshaledScript.UpdatedAt == scriptRecord.UpdatedAt

	if !isScriptEqual {
		t.Fatalf("InstallerScript roundtrip mismatch: got %+v, want %+v", unmarshaledScript, scriptRecord)
	}

	versionRecord := InstallerVersion{
		ID:           10,
		ScriptID:     1,
		Slug:         "nodejs-installer",
		Version:      "1.0.0",
		TargetOS:     "win",
		Instructions: `{"steps":["download","install"]}`,
		CreatedAt:    "2026-08-26T00:00:00Z",
	}

	marshaledVersion, versionErr := json.Marshal(versionRecord)
	if versionErr != nil {
		t.Fatalf("json.Marshal(InstallerVersion) failed: %v", versionErr)
	}

	var unmarshaledVersion InstallerVersion
	unmarshalVersionErr := json.Unmarshal(marshaledVersion, &unmarshaledVersion)
	if unmarshalVersionErr != nil {
		t.Fatalf("json.Unmarshal(InstallerVersion) failed: %v", unmarshalVersionErr)
	}

	isVersionEqual := unmarshaledVersion.ID == versionRecord.ID &&
		unmarshaledVersion.ScriptID == versionRecord.ScriptID &&
		unmarshaledVersion.Slug == versionRecord.Slug &&
		unmarshaledVersion.Version == versionRecord.Version &&
		unmarshaledVersion.TargetOS == versionRecord.TargetOS &&
		unmarshaledVersion.Instructions == versionRecord.Instructions &&
		unmarshaledVersion.CreatedAt == versionRecord.CreatedAt

	if !isVersionEqual {
		t.Fatalf("InstallerVersion roundtrip mismatch: got %+v, want %+v", unmarshaledVersion, versionRecord)
	}
}

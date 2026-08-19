package cluster

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestPreflightBoxRender(t *testing.T) {
	effective := []ClusterNode{
		{ID: "node1", DisplayId: 1, IP: "192.168.1.1"},
		{ID: "node2", DisplayId: 2, IP: "192.168.1.2"},
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run preflight with autoConfirm=true to avoid blocking
	PrintPreflight(ServersClients, effective, `ps "echo hello"`, "RUN-20260819-001", true)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !bytes.Contains([]byte(output), []byte("Run Ref:")) {
		t.Errorf("Missing Run Ref in output: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("RUN-20260819-001")) {
		t.Errorf("Missing Run Ref value in output: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("Command:")) {
		t.Errorf("Missing Command in output: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte(`ps "echo hello"`)) {
		t.Errorf("Missing Command value in output: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("Selector:")) {
		t.Errorf("Missing Selector in output: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("Servers & Clients")) {
		t.Errorf("Missing Selector value in output: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("Nodes:")) {
		t.Errorf("Missing Nodes in output: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("2 effective targets")) {
		t.Errorf("Missing Nodes value in output: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("node1 (ID: 1, IP: 192.168.1.1)")) {
		t.Errorf("Missing node1 in output: %s", output)
	}
}

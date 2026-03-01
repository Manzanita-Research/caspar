package output

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
)

func captureStdout(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestJSON(t *testing.T) {
	out := captureStdout(func() {
		JSON(map[string]string{"key": "value"})
	})

	var result map[string]string
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if result["key"] != "value" {
		t.Errorf("expected key=value, got %v", result)
	}
}

func TestPrint_JSONMode(t *testing.T) {
	out := captureStdout(func() {
		Print(true, map[string]string{"mode": "json"}, func() {
			t.Error("human function should not be called in JSON mode")
		})
	})

	if !strings.Contains(out, `"mode"`) {
		t.Errorf("expected JSON output, got: %s", out)
	}
}

func TestPrint_HumanMode(t *testing.T) {
	called := false
	captureStdout(func() {
		Print(false, nil, func() {
			called = true
		})
	})

	if !called {
		t.Error("expected human function to be called")
	}
}

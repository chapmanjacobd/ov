package oviewer

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestIsTermux(t *testing.T) {
	tests := []struct {
		name   string
		env    map[string]string
		expect bool
	}{
		{
			name:   "linux",
			env:    map[string]string{"TERMUX_VERSION": "", "PREFIX": "/usr"},
			expect: false,
		},
		{
			name:   "termux version",
			env:    map[string]string{"TERMUX_VERSION": "0.118.3", "PREFIX": ""},
			expect: true,
		},
		{
			name:   "termux prefix",
			env:    map[string]string{"TERMUX_VERSION": "", "PREFIX": "/data/data/com.termux/files/usr"},
			expect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.env {
				t.Setenv(key, value)
			}
			if got := isTermux(); got != tt.expect {
				t.Errorf("isTermux() = %v, want %v", got, tt.expect)
			}
		})
	}
}

func TestTermuxClipboard(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("fake clipboard commands are shell scripts")
	}

	dir := t.TempDir()
	clipboardFile := filepath.Join(dir, "clipboard.txt")
	t.Setenv("TERMUX_CLIPBOARD_TEST_FILE", clipboardFile)
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))

	writeScript(t, filepath.Join(dir, "termux-clipboard-get"), "#!/bin/sh\ncat \"$TERMUX_CLIPBOARD_TEST_FILE\"\n")
	writeScript(t, filepath.Join(dir, "termux-clipboard-set"), "#!/bin/sh\ncat > \"$TERMUX_CLIPBOARD_TEST_FILE\"\n")

	ctx := context.Background()
	if err := initTermuxClipboard(); err != nil {
		t.Fatalf("initTermuxClipboard() error = %v", err)
	}
	if err := writeTermuxClipboard(ctx, "copied text"); err != nil {
		t.Fatalf("writeTermuxClipboard() error = %v", err)
	}
	got, err := readTermuxClipboard(ctx)
	if err != nil {
		t.Fatalf("readTermuxClipboard() error = %v", err)
	}
	if string(got) != "copied text" {
		t.Errorf("readTermuxClipboard() = %q, want %q", got, "copied text")
	}
}

func TestInitTermuxClipboardMissing(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	err := initTermuxClipboard()
	if !errors.Is(err, errTermuxAPIMissing) {
		t.Errorf("initTermuxClipboard() error = %v, want %v", err, errTermuxAPIMissing)
	}
}

func writeScript(t *testing.T, path, content string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatalf("write script: %v", err)
	}
}

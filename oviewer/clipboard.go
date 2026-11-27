package oviewer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// errTermuxAPIMissing explains how to enable the system clipboard in Termux.
var errTermuxAPIMissing = errors.New("clipboard: install termux-api (pkg install termux-api) to use the system clipboard")

// Clipboard access is split into a Termux backend and a native backend.
// Termux exposes the Android clipboard through the termux-api commands.
// The native backend's Android implementation requires gomobile and does not
// build with the Termux Go toolchain, so it is omitted there (see
// clipboard_android.go).

// isTermux reports whether ov is running inside Termux.
func isTermux() bool {
	return os.Getenv("TERMUX_VERSION") != "" ||
		strings.Contains(os.Getenv("PREFIX"), "com.termux")
}

// initClipboard prepares the clipboard backend.
func initClipboard() error {
	if isTermux() {
		return initTermuxClipboard()
	}
	return initNativeClipboard()
}

// readClipboard returns the text currently on the clipboard.
func readClipboard(ctx context.Context) ([]byte, error) {
	if isTermux() {
		return readTermuxClipboard(ctx)
	}
	return readNativeClipboard(ctx)
}

// writeClipboard replaces the text on the clipboard.
func writeClipboard(ctx context.Context, text string) error {
	if isTermux() {
		return writeTermuxClipboard(ctx, text)
	}
	return writeNativeClipboard(ctx, text)
}

// initTermuxClipboard checks that the termux-api clipboard commands exist.
func initTermuxClipboard() error {
	for _, name := range []string{"termux-clipboard-get", "termux-clipboard-set"} {
		if _, err := exec.LookPath(name); err != nil {
			return fmt.Errorf("%w: %v", errTermuxAPIMissing, err)
		}
	}
	return nil
}

func readTermuxClipboard(ctx context.Context) ([]byte, error) {
	out, err := exec.CommandContext(ctx, "termux-clipboard-get").Output()
	if err != nil {
		return nil, fmt.Errorf("clipboard: termux-clipboard-get: %w", err)
	}
	return out, nil
}

func writeTermuxClipboard(ctx context.Context, text string) error {
	cmd := exec.CommandContext(ctx, "termux-clipboard-set")
	cmd.Stdin = strings.NewReader(text)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clipboard: termux-clipboard-set: %w", err)
	}
	return nil
}

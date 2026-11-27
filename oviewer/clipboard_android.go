//go:build android

package oviewer

import (
	"context"
	"errors"
)

// errNativeAndroidClipboard is returned outside Termux, where a command-line
// process has no native Android clipboard without a gomobile app.
var errNativeAndroidClipboard = errors.New("clipboard: native Android clipboard requires a gomobile app; run ov in Termux with termux-api installed")

func initNativeClipboard() error {
	return nil
}

func readNativeClipboard(context.Context) ([]byte, error) {
	return nil, errNativeAndroidClipboard
}

func writeNativeClipboard(context.Context, string) error {
	return errNativeAndroidClipboard
}

//go:build !android

package oviewer

import (
	"context"

	"golang.design/x/clipboard"
)

func initNativeClipboard() error {
	return clipboard.Init()
}

func readNativeClipboard(ctx context.Context) ([]byte, error) {
	return clipboard.Read(ctx, clipboard.FmtText)
}

func writeNativeClipboard(ctx context.Context, text string) error {
	_, err := clipboard.Write(ctx, clipboard.FmtText, []byte(text))
	return err
}

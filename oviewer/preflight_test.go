package oviewer

import (
	"context"
	"io"
	"path/filepath"
	"testing"
	"time"

	"github.com/gdamore/tcell/v3"
)

func TestRoot_quitSmallPreflightFilterFits(t *testing.T) {
	tcellNewScreen = fakeScreen
	defer func() {
		tcellNewScreen = tcell.NewScreen
	}()

	root := rootFileReadHelper(t, filepath.Join(testdata, "test.txt"))
	root.Config.QuitSmall = true
	root.Config.QuitSmallFilter = true
	if err := root.prepareRun(context.Background()); err != nil {
		t.Fatal(err)
	}
	root.prepareScreen()
	root.Filter("test", false)

	if !root.quitSmallPreflight(context.Background()) {
		t.Fatalf("quitSmallPreflight() = false, want true (docNum=%d type=%v eof=%v lines=%d)",
			root.DocumentLen(), root.Doc.documentType, root.Doc.BufEOF(), root.Doc.BufEndNum())
	}
	if !root.Config.IsWriteOnExit {
		t.Error("IsWriteOnExit = false, want true")
	}
	if len(root.OnExit) == 0 {
		t.Fatal("OnExit is empty")
	}
}

func TestRoot_quitSmallPreflightFilterDoesNotFit(t *testing.T) {
	tcellNewScreen = fakeScreen
	defer func() {
		tcellNewScreen = tcell.NewScreen
	}()

	root := rootFileReadHelper(t, filepath.Join(testdata, "test3.txt"))
	root.Config.QuitSmall = true
	root.Config.QuitSmallFilter = true
	if err := root.prepareRun(context.Background()); err != nil {
		t.Fatal(err)
	}
	root.prepareScreen()
	root.Filter("1", false)

	if root.quitSmallPreflight(context.Background()) {
		t.Fatal("quitSmallPreflight() = true, want false")
	}
	if root.Config.IsWriteOnExit {
		t.Error("IsWriteOnExit = true, want false")
	}
}

func TestRoot_quitSmallPreflightWaitsForEOF(t *testing.T) {
	tcellNewScreen = fakeScreen
	defer func() {
		tcellNewScreen = tcell.NewScreen
	}()

	pr, pw := io.Pipe()
	doc, err := NewDocument()
	if err != nil {
		t.Fatal(err)
	}
	if err := doc.ControlReader(pr, nil); err != nil {
		t.Fatal(err)
	}
	root, err := NewOviewer(doc)
	if err != nil {
		t.Fatal(err)
	}
	setFixedTestScreen(t, root)
	root.Config.QuitSmall = true
	if err := root.prepareRun(context.Background()); err != nil {
		t.Fatal(err)
	}
	root.prepareScreen()

	go func() {
		if _, err := pw.Write([]byte("a\nb\n")); err != nil {
			t.Error(err)
		}
		time.Sleep(20 * time.Millisecond)
		pw.Close()
	}()

	if !root.quitSmallPreflight(context.Background()) {
		t.Fatal("quitSmallPreflight() = false, want true after EOF")
	}
	if !root.Config.IsWriteOnExit {
		t.Error("IsWriteOnExit = false, want true")
	}
	if len(root.OnExit) != 2 {
		t.Errorf("len(OnExit) = %d, want 2", len(root.OnExit))
	}
}

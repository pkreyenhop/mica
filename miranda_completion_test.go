package main

import (
	"strings"
	"testing"

	"mica/editor"
)

func TestTryManualCompletionMirandaKeyword(t *testing.T) {
	app := &appState{}
	app.initBuffers(editor.NewEditor("oth"))
	app.currentPath = "test.m"
	app.buffers[0].path = app.currentPath
	app.ed.Caret = len([]rune("oth"))

	if !tryManualCompletion(app) {
		t.Fatal("expected completion to apply")
	}
	if got := app.ed.String(); got != "otherwise" {
		t.Fatalf("completed text = %q, want %q", got, "otherwise")
	}
}

func TestTryManualCompletionCurrentFileDefinition(t *testing.T) {
	src := "myvalue x = x\nmain = myv"
	app := &appState{}
	app.initBuffers(editor.NewEditor(src))
	app.currentPath = "test.m"
	app.buffers[0].path = app.currentPath
	app.ed.Caret = len([]rune(src))

	if !tryManualCompletion(app) {
		t.Fatal("expected completion to apply")
	}
	if got := app.ed.String(); !strings.Contains(got, "main = myvalue") {
		t.Fatalf("completed text = %q, want current-file symbol expansion", got)
	}
}

func TestTryManualCompletionStdlibDefinition(t *testing.T) {
	app := &appState{}
	app.initBuffers(editor.NewEditor("dropw"))
	app.currentPath = "test.m"
	app.buffers[0].path = app.currentPath
	app.ed.Caret = len([]rune("dropw"))

	if !tryManualCompletion(app) {
		t.Fatal("expected completion to apply")
	}
	if got := app.ed.String(); got != "dropwhile" {
		t.Fatalf("completed text = %q, want %q", got, "dropwhile")
	}
}

func TestTryManualCompletionAmbiguousShowsPopup(t *testing.T) {
	app := &appState{}
	app.initBuffers(editor.NewEditor("f"))
	app.currentPath = "test.m"
	app.buffers[0].path = app.currentPath
	app.ed.Caret = len([]rune("f"))

	if !tryManualCompletion(app) {
		t.Fatal("expected completion handling")
	}
	if !app.completionPopup.active {
		t.Fatal("expected completion popup for ambiguous prefix")
	}
	if len(app.completionPopup.items) < 2 {
		t.Fatalf("expected multiple completion items, got %d", len(app.completionPopup.items))
	}
}

func TestTryManualCompletionNonMirandaDoesNothing(t *testing.T) {
	app := &appState{}
	app.initBuffers(editor.NewEditor("oth"))
	app.currentPath = "note.md"
	app.buffers[0].path = app.currentPath
	app.ed.Caret = len([]rune("oth"))

	if tryManualCompletion(app) {
		t.Fatal("completion should be inactive outside miranda mode")
	}
	if got := app.ed.String(); got != "oth" {
		t.Fatalf("buffer changed unexpectedly: %q", got)
	}
}

func TestMirandaCompletionItemsPreferCurrentFileOverStdlib(t *testing.T) {
	src := "map x = x\nmain = ma"
	app := &appState{}
	app.initBuffers(editor.NewEditor(src))
	app.currentPath = "test.m"
	app.buffers[0].path = app.currentPath
	items := mirandaCompletionItems(app, []rune(src), "ma")
	if len(items) == 0 {
		t.Fatal("expected completion items")
	}
	for _, it := range items {
		if it.Label == "map" {
			if it.Detail != "current file" {
				t.Fatalf("map detail = %q, want current file", it.Detail)
			}
			return
		}
	}
	t.Fatal("missing map completion")
}

func TestLongestSharedPrefix(t *testing.T) {
	items := []completionItem{
		{Insert: "drop"},
		{Insert: "dropwhile"},
		{Insert: "dropdown"},
	}
	if got := longestSharedPrefix(items); got != "drop" {
		t.Fatalf("longestSharedPrefix = %q, want drop", got)
	}
}

func TestMirandaLocalDefinitionsCacheInvalidatesOnTextRev(t *testing.T) {
	app := &appState{}
	app.initBuffers(editor.NewEditor("foo x = x\n"))
	app.currentPath = "test.m"
	app.buffers[0].path = app.currentPath

	first := mirandaLocalDefinitions(app, app.ed.Runes())
	if _, ok := first["foo"]; !ok {
		t.Fatal("expected foo definition")
	}

	app.ed.SetRunes([]rune("bar x = x\n"))
	app.touchActiveBufferText()
	second := mirandaLocalDefinitions(app, app.ed.Runes())
	if _, ok := second["bar"]; !ok {
		t.Fatal("expected bar definition after text revision change")
	}
	if _, ok := second["foo"]; ok {
		t.Fatal("stale foo definition should not persist after cache refresh")
	}
}

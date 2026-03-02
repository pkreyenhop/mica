package main

import (
	"strings"
	"testing"

	"mica/editor"
)

func TestShowSymbolInfoShowsCurrentBufferDefinition(t *testing.T) {
	src := "foo x = x+1\nbar = foo 2\n"
	app := appState{}
	app.initBuffers(editor.NewEditor(src))
	app.buffers[0].mode = syntaxMiranda

	caret := strings.Index(src, "foo 2")
	if caret < 0 {
		t.Fatal("test setup failed to find symbol")
	}
	app.ed.Caret = caret

	got := showSymbolInfo(&app)
	if !strings.Contains(got, "Definition (current buffer)") {
		t.Fatalf("expected current-buffer definition header, got:\n%s", got)
	}
	if !strings.Contains(got, "foo x = x+1") {
		t.Fatalf("expected function definition code, got:\n%s", got)
	}
}

func TestShowSymbolInfoShowsStdlibDefinition(t *testing.T) {
	src := "map"
	app := appState{}
	app.initBuffers(editor.NewEditor(src))
	app.buffers[0].mode = syntaxMiranda
	app.ed.Caret = 1

	got := showSymbolInfo(&app)
	if !strings.Contains(got, "miranda/miralib/stdenv.m") {
		t.Fatalf("expected stdenv source path in popup, got:\n%s", got)
	}
	if !strings.Contains(strings.ToLower(got), "map") {
		t.Fatalf("expected map definition in popup, got:\n%s", got)
	}
}

func TestShowSymbolInfoShowsStdlibDefinitionReverse(t *testing.T) {
	src := "reverse"
	app := appState{}
	app.initBuffers(editor.NewEditor(src))
	app.buffers[0].mode = syntaxMiranda
	app.ed.Caret = 1

	got := showSymbolInfo(&app)
	if !strings.Contains(got, "miranda/miralib/stdenv.m") {
		t.Fatalf("expected stdenv source path in popup, got:\n%s", got)
	}
	if !strings.Contains(strings.ToLower(got), "reverse") {
		t.Fatalf("expected reverse definition in popup, got:\n%s", got)
	}
}

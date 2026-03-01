package main

import (
	"testing"

	"mica/editor"
)

func TestDetectSyntaxByPath(t *testing.T) {
	tests := []struct {
		path string
		want syntaxKind
	}{
		{path: "a.md", want: syntaxMarkdown},
		{path: "a.markdown", want: syntaxMarkdown},
		{path: "a.m", want: syntaxMiranda},
		{path: "a.go", want: syntaxNone},
		{path: "a.c", want: syntaxNone},
		{path: "a.h", want: syntaxNone},
		{path: "a.txt", want: syntaxNone},
	}
	for _, tc := range tests {
		if got := detectSyntax(tc.path, ""); got != tc.want {
			t.Fatalf("detectSyntax(%q)=%v, want %v", tc.path, got, tc.want)
		}
	}
}

func TestDetectSyntaxByContent(t *testing.T) {
	if got := detectSyntax("", "## title\ntext\n"); got != syntaxMarkdown {
		t.Fatalf("markdown heading detectSyntax=%v, want %v", got, syntaxMarkdown)
	}
	if got := detectSyntax("", "plain text\n"); got != syntaxNone {
		t.Fatalf("plain text detectSyntax=%v, want %v", got, syntaxNone)
	}
}

func TestSyntaxKindLabel(t *testing.T) {
	tests := []struct {
		kind syntaxKind
		want string
	}{
		{kind: syntaxNone, want: "text"},
		{kind: syntaxMarkdown, want: "markdown"},
		{kind: syntaxMiranda, want: "miranda"},
	}
	for _, tc := range tests {
		if got := syntaxKindLabel(tc.kind); got != tc.want {
			t.Fatalf("syntaxKindLabel(%v)=%q, want %q", tc.kind, got, tc.want)
		}
	}
}

func TestBufferSyntaxKindUsesForcedMode(t *testing.T) {
	app := appState{}
	app.initBuffers(editor.NewEditor("plain text"))
	app.currentPath = "note.txt"
	app.buffers[0].path = "note.txt"

	if got := bufferSyntaxKind(&app, app.currentPath, app.ed.Runes()); got != syntaxNone {
		t.Fatalf("default syntax kind=%v, want text/none", got)
	}
	app.buffers[0].mode = syntaxMarkdown
	if got := bufferSyntaxKind(&app, app.currentPath, app.ed.Runes()); got != syntaxMarkdown {
		t.Fatalf("forced syntax kind=%v, want markdown", got)
	}
}

func TestCycleBufferModeSkipsGoAndC(t *testing.T) {
	app := appState{}
	app.initBuffers(editor.NewEditor(""))

	if got := cycleBufferMode(&app); got != "markdown" {
		t.Fatalf("first cycle=%q, want markdown", got)
	}
	if app.buffers[0].mode != syntaxMarkdown {
		t.Fatalf("mode after first cycle=%v, want %v", app.buffers[0].mode, syntaxMarkdown)
	}
	if got := cycleBufferMode(&app); got != "miranda" {
		t.Fatalf("second cycle=%q, want miranda", got)
	}
	if app.buffers[0].mode != syntaxMiranda {
		t.Fatalf("mode after second cycle=%v, want %v", app.buffers[0].mode, syntaxMiranda)
	}
	if got := cycleBufferMode(&app); got != "text" {
		t.Fatalf("third cycle=%q, want text", got)
	}
	if app.buffers[0].mode != syntaxNone {
		t.Fatalf("mode after third cycle=%v, want %v", app.buffers[0].mode, syntaxNone)
	}
}

package main

import (
	"strings"
	"testing"

	"mica/editor"

	"github.com/gdamore/tcell/v2"
)

func TestCtrlRuneToKey(t *testing.T) {
	if got, ok := ctrlRuneToKey('q'); !ok || got != keyQ {
		t.Fatalf("ctrlRuneToKey('q') = %v %v, want keyQ true", got, ok)
	}
	if got, ok := ctrlRuneToKey('i'); !ok || got != keyTab {
		t.Fatalf("ctrlRuneToKey('i') = %v %v, want keyTab true", got, ok)
	}
}

func TestEscHelpPopupShowsOnlyEscCommands(t *testing.T) {
	lines := escHelpPopupLines()
	if len(lines) == 0 {
		t.Fatal("expected non-empty Esc help lines")
	}
	all := strings.Join(lines, "\n")
	if strings.Contains(all, "Ctrl+") || strings.Contains(all, "ctrl+") {
		t.Fatalf("Esc help should not list Ctrl sequences, got:\n%s", all)
	}
	if !strings.Contains(all, "m  cycle language mode") {
		t.Fatalf("Esc help should list next-letter commands, got:\n%s", all)
	}
}

func TestHandleTUIEscPrefixCreatesBuffer(t *testing.T) {
	app := appState{}
	app.initBuffers(editor.NewEditor(""))

	if !handleTUIKey(&app, tcell.NewEventKey(tcell.KeyEscape, 0, 0)) {
		t.Fatal("Esc should arm prefix")
	}
	if !handleTUIKey(&app, tcell.NewEventKey(tcell.KeyRune, 'b', 0)) {
		t.Fatal("Esc+b should continue")
	}
	if len(app.buffers) != 2 {
		t.Fatalf("Esc+b should create buffer, got %d", len(app.buffers))
	}
}

func TestRenderDataTreatsPathAsMirandaWhenUnforced(t *testing.T) {
	app := appState{syntaxHL: newMirandaHighlighter()}
	app.initBuffers(editor.NewEditor("package main\n"))
	app.currentPath = "p.go"
	app.buffers[0].path = "p.go"

	_, _, lang, _ := renderData(&app)
	if lang != "miranda" {
		t.Fatalf("lang=%q, want miranda", lang)
	}
}

func TestRenderDataRespectsForcedMarkdownMode(t *testing.T) {
	app := appState{syntaxHL: newMirandaHighlighter()}
	app.initBuffers(editor.NewEditor("plain text\n"))
	app.currentPath = "notes.txt"
	app.buffers[0].path = "notes.txt"
	app.buffers[0].mode = syntaxMarkdown

	_, _, lang, _ := renderData(&app)
	if lang != "markdown" {
		t.Fatalf("lang=%q, want markdown", lang)
	}
}

func TestDrawTUIDoesNotPanic(t *testing.T) {
	s := tcell.NewSimulationScreen("UTF-8")
	if err := s.Init(); err != nil {
		t.Fatalf("init simulation screen: %v", err)
	}
	defer s.Fini()
	s.SetSize(80, 24)

	app := appState{syntaxHL: newMirandaHighlighter()}
	app.initBuffers(editor.NewEditor("hello\nworld\n"))
	drawTUI(s, &app)
}

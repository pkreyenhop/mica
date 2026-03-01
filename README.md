# Mica

Mica (Miranda Cat) is a terminal-first text editor with Canon Cat inspired navigation, fast buffer workflows, and a headless editing core.

## Build and Run

```sh
go build -o mica .
./mica [file1 file2 ...]
```

- Missing startup files open as empty buffers and are created on first save.
- `Shift+Tab` cycles buffers.

## Core Shortcuts

- `Ctrl+B`: new buffer
- `Ctrl+O`: open file picker buffer
- `Ctrl+L`: open file/directory under caret in picker
- `Esc+W`: write/save-as prompt
- `Esc+Shift+S`: save dirty buffers
- `Ctrl+Q`: close current buffer
- `Esc+Shift+Q`: quit all
- `Esc+Esc`: close current buffer (prefix mode)
- `Esc+M`: cycle forced language mode (`text -> markdown -> miranda -> text`)
- `Esc+I`: Miranda symbol info popup (shows function definitions from current buffer or vendored stdlib when available)
- `Esc+/`: incremental search (`/` lock, `Tab` next, `Shift+Tab` previous)
- `Esc+X`: line-highlight mode (`x` extends)
- `Esc+Space`: less mode (`Space` page, `Esc` exit)
- `Esc+Shift+Delete`: clear active buffer contents
- `Ctrl+/`: toggle line comments on selection/current line
- `Ctrl+A`/`Ctrl+E`: line start/end
- `Ctrl+Shift+A`/`Ctrl+Shift+E`: buffer start/end
- `Ctrl+,` / `Ctrl+.`: page up/down
- `Ctrl+K`: kill to end of line
- `Ctrl+U`: undo
- `Ctrl+C` / `Ctrl+X` / `Ctrl+V`: copy/cut/paste

## Syntax Modes

- Auto-detect supports Markdown (`.md`, `.markdown`) and Miranda (`.m`).
- Plain text is the default.
- `Esc+M` sets a forced mode per buffer.

## Project Layout

- `main_tui.go`: tcell frontend loop and terminal rendering.
- `main.go`: shared app state and non-UI helpers.
- `input_core.go`: platform-neutral key/text command handling.
- `editor/`: headless editing core (gap-buffer-backed storage).
- `miranda/`: vendored Miranda runtime/source tree tracked in this repository.

## Tests

```sh
go test ./editor
go test ./...
```

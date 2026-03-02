# Mica

Mica (Miranda Cat) is a terminal-first text editor with Canon Cat inspired navigation, fast buffer workflows, and a headless editing core.

## Build and Run

```sh
go build -o mica .
./mica [file1 file2 ...]
./mica +13 test.m
```

- Missing startup files open as empty buffers and are created on first save.
- `+<line>` before a file path opens that file and places the caret on the 1-based line number.
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
- `Esc+M`: cycle forced language mode (`auto/miranda -> markdown -> miranda -> auto/miranda`)
- `Esc+I`: Miranda symbol info popup (shows function definitions from current buffer or vendored stdlib when available)
- `Tab`: Miranda completion for keywords, builtins, current-file symbols, and vendored stdlib symbols (opens chooser when ambiguous)
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
- Miranda is the default mode for non-Markdown buffers.
- `Esc+M` sets a forced mode per buffer.

## Project Layout

- `main_tui.go`: tcell frontend loop and terminal rendering.
- `main.go`: shared app state and non-UI helpers.
- `input_core.go`: platform-neutral key/text command handling.
- `miranda_highlight_common.go` + `miranda_highlight.go`: Miranda/Markdown syntax detection and lexical highlighting.
- `miranda_symbol_info.go`: `Esc+I` symbol lookup and stdlib help extraction.
- `miranda_completion.go` + `miranda_completion_item.go`: Miranda completion candidate assembly and popup item payload type.
- `editor/`: headless editing core (gap-buffer-backed storage).
- `miranda/`: vendored Miranda runtime/source tree tracked in this repository.

## Tests

```sh
go test ./editor
go test ./...
```

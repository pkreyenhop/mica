# Mica Manual

## Overview

Mica is a TUI editor with a headless core. It focuses on keyboard-first editing, multi-buffer flow, and lightweight command-prefix control.

## Launch

```sh
go build -o mica .
./mica [file1 file2 ...]
```

## Editing and Movement

- `Backspace`: delete backward
- `Delete`: delete word under/left of caret
- `Shift+Delete`: delete current line
- `Enter`: newline
- `Ctrl+K`: kill to end of line
- `Ctrl+U`: undo
- `Ctrl+/`: toggle comments (`//`) on selection/current line
- `Ctrl+A` / `Ctrl+E`: line start/end
- `Ctrl+Shift+A` / `Ctrl+Shift+E`: buffer start/end
- `Ctrl+,` / `Ctrl+.`: page up/down
- Arrows/PageUp/PageDown: caret navigation (Shift extends selection)

## Buffers and Files

- `Ctrl+B`: new buffer
- `Shift+Tab`: cycle buffers
- `Ctrl+O`: open picker buffer rooted at current root
- `Ctrl+L`: open selected picker entry
- `Esc+W`: save-as prompt
- `Esc+Shift+S`: save dirty buffers
- `Ctrl+Q`: close buffer
- `Esc+Shift+Q`: quit all
- Missing startup file paths become empty buffers and are created on save.

## Esc Prefix

- `Esc` arms prefix mode.
- `Esc+Esc`: close current buffer
- `Esc+W`: write prompt
- `Esc+Shift+S`: save dirty buffers
- `Esc+Shift+Q`: quit all
- `Esc+M`: cycle forced language mode (`text -> markdown -> miranda -> text`)
- `Esc+/`: search mode
- `Esc+X`: line-highlight mode
- `Esc+Space`: less mode
- `Esc+Shift+Delete`: clear active buffer

If `Esc` remains pending briefly, a lower-right helper popup shows available next-letter commands.

## Search and Highlight Modes

- Search: `Esc+/`, type query, `/` locks pattern, `Tab`/`Shift+Tab` navigate matches.
- Line highlight: `Esc+X`, then `x` extends selection by line.
- `Esc` exits either mode.

## Syntax Modes

- Auto detect: Markdown (`.md`, `.markdown`), Miranda (`.m`), else text.
- Forced mode cycle with `Esc+M`.

## Rendering

- Status line shows buffer info, mode label, cwd, and dirty marker.
- Bottom line shows prompt/search/auxiliary mode text.
- Line numbers render in a gutter.

## Testing

```sh
go test ./editor
go test ./...
```

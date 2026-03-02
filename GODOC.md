# GoDoc

Generated with:

```sh
go doc -cmd .
go doc -all ./editor
```

## `mica` (command package)

```text
package main // import "."
```

### Runtime Modes (TUI)

- `Esc` command-prefix mode with delayed lower-right helper popup.
- `Esc+I` symbol info popup for Miranda symbols and function definitions.
- `Tab` completion in Miranda buffers for keywords, builtins, current-file symbols, and vendored stdlib symbols.
- `Esc+/` incremental search mode (`/` lock, `Tab`/`Shift+Tab` navigation).
- `Esc+X` line-highlight mode (`x` extends selection, `Esc` exits).
- `Esc+Space` less mode (`Space` page, `Esc` exit).
- Startup CLI accepts `+<line>` before a file path to place caret at that 1-based line.

### Syntax Highlighting

- Syntax highlighting uses built-in lexical tokenizers (no LSP/Tree-sitter dependency).
- Active syntax kinds: miranda (default), markdown.
- Implementation files: `miranda_highlight_common.go`, `miranda_highlight.go`.

## `mica/editor`

```text
package editor // import "."

Package editor provides headless editing and Canon-Cat-inspired Leap logic.
```

Use `go doc -all ./editor` for the current API surface.

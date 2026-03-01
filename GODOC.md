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
- `Esc+/` incremental search mode (`/` lock, `Tab`/`Shift+Tab` navigation).
- `Esc+X` line-highlight mode (`x` extends selection, `Esc` exits).
- `Esc+Space` less mode (`Space` page, `Esc` exit).

### Syntax Highlighting

- Syntax highlighting uses pure-Go Tree-sitter (`github.com/odvcencio/gotreesitter`).
- Active syntax kinds: text (default), markdown, miranda.

## `mica/editor`

```text
package editor // import "."

Package editor provides headless editing and Canon-Cat-inspired Leap logic.
```

Use `go doc -all ./editor` for the current API surface.

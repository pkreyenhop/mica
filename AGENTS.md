# Repository Guidelines

## Project Structure & Module Organization
- Identity: the editor is called `mica` (Miranda Cat) and draws inspiration from the Canon Cat, Helix, acme, AMP, and Emacs.
- `main_tui.go` — TUI frontend (tcell): screen setup, event loop, rendering, and key dispatch.
- `main.go` — shared app state and non-UI helpers.
- `input_core.go` — platform-agnostic input/controller layer.
- `editor/` — UI-free core: buffer management, leap/search, selection, clipboard abstraction.
- Root tests: `main_open_test.go`, `main_buffer_test.go`, `main_scroll_test.go`, `main_syntax_test.go`, `main_tui_test.go`, `main_help_test.go`.
- `RULES.md` — canonical behavior list.
- `README.md` — user-facing overview.
- `miranda/` — vendored Miranda runtime/source tree; treat as part of `mica` (not a submodule).

## Build, Test, and Development Commands
- `go run .` — launch the TUI editor.
- `go build .` — compile the editor.
- `go test ./editor` — run headless logic tests.
- `go test ./...` — run full tests.

## Coding Style & Naming Conventions
- Run `gofmt` before sending changes.
- Keep core logic UI-agnostic inside `editor/`.
- Naming: directions as `DirFwd`/`DirBack`; caret/selection fields as `Caret`, `Sel`, `Leap`.
- Use `[]rune` for buffer text.
- Prefer editor accessors (`Runes`, `String`, `RuneLen`, `SetRunes`) over direct storage mutation.
- Keep `helpEntries`, README, and RULES aligned.
- Esc-prefixed command shortcuts include `Esc+I` for Miranda symbol-definition popup, `Esc+M` for mode cycle (`text/markdown/miranda`), and `Esc+Shift+Delete` for buffer clear.

## Testing Guidelines
- Prefer scenario-style behavior tests.
- Keep tests deterministic and headless.
- Add key-flow checks in `main_tui_test.go` and command-mode checks in `input_core_test.go`.

## Commit & Pull Request Guidelines
- Commit messages should be imperative and present tense.
- Mention key behaviors changed and commands run (`go test ./...`).

## Environment & Configuration Notes
- Go 1.26+ per `go.mod`.

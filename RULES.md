# Mica Behaviour Rules

Mica (Miranda Cat) is a keyboard-centric editor inspired by the Canon Cat, Helix, acme, AMP, and Emacs.

- **Leap navigation**
  - Leap trigger keys are currently unbound in TUI mode.
  - Leap selection/repeat logic remains in the editor core.

- **Buffers and files**
  - `Ctrl+B` creates `<untitled>` buffers.
  - `Shift+Tab` cycles buffers.
  - `Ctrl+O` opens a picker buffer rooted at current root.
  - `Ctrl+L` opens the picker line path (or switches to an already-open buffer).
  - Startup accepts multiple file paths; missing paths open empty buffers and are created on save.
  - Startup accepts `+<line>` before a file path (for example `mica +13 test.m`) and moves caret to that 1-based line.
  - `Esc+W` opens save-as prompt.
  - `Esc+Shift+S` saves dirty buffers.
  - `Ctrl+Q` closes current buffer.
  - `Esc+Shift+Q` quits all.

- **Esc command prefix**
  - `Esc` arms command mode.
  - `Esc+Esc` closes current buffer.
  - `Esc+I` opens Miranda symbol info popup and shows symbol/function definitions when found.
  - `Esc+M` cycles forced language mode: `auto/miranda -> markdown -> miranda -> auto/miranda`.
  - `Esc+/` enters incremental search.
  - `Esc+X` enters line-highlight mode.
  - `Esc+Space` enters less mode.
  - `Esc+Shift+Delete` clears active buffer contents and marks dirty.
  - If no second key follows quickly, a lower-right helper popup appears.

- **Editing and movement**
  - Text input inserts runes.
  - Double-space inserts a tab at line start.
  - `Tab` in Miranda buffers expands prefixes using keywords/builtins/current-file symbols/stdlib symbols.
  - Ambiguous Miranda `Tab` completion opens a chooser popup; `Tab`/`Shift+Tab` cycle and `Enter` applies.
  - `Backspace` deletes backward.
  - `Delete` removes word under/left of caret.
  - `Shift+Delete` removes current line.
  - `Ctrl+A`/`Ctrl+E` move to line start/end.
  - `Ctrl+Shift+A`/`Ctrl+Shift+E` move to buffer start/end.
  - `Ctrl+,` / `Ctrl+.` page up/down.
  - `Ctrl+K` kills to end of line.
  - `Ctrl+U` performs single-step undo.
  - `Ctrl+/` toggles `//` comments on selection/current line.
  - `Ctrl+C`/`Ctrl+X`/`Ctrl+V` perform copy/cut/paste.

- **Search and highlight**
  - In search mode, `/` locks query and `Tab`/`Shift+Tab` navigate matches with wrap.
  - Locking `/` on empty query repeats last non-empty query.
  - In locked search mode, `x` switches to line-highlight mode.
  - In line-highlight mode, repeated `x` extends selection by one line.

- **Syntax and rendering**
  - Auto detection supports Markdown (`.md`, `.markdown`); otherwise Miranda.
  - Status bar shows mode/lang/cwd/dirty marker.
  - Input line is used for prompts and mode text.
  - Gutter shows line numbers with current-line emphasis.

- **Dirty tracking**
  - Editing operations mark the active buffer dirty.
  - Save/load clears dirty.

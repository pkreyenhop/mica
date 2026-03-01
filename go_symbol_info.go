package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"
)

var mirandaKeywordInfo = map[string]string{
	"where":     "introduces local definitions for a declaration or expression",
	"if":        "conditional expression guard",
	"otherwise": "default guard branch",
	"type":      "declares a type synonym or algebraic type",
	"abstype":   "declares an abstract data type",
	"with":      "introduces constructor/export detail in abstype declarations",
	"div":       "integer division operator",
	"mod":       "integer remainder operator",
	"include":   "compiler/include directive in scripts",
	"export":    "compiler/export directive in scripts",
	"free":      "compiler/free directive",
	"show":      "show/printing related keyword in parser tokens",
	"readvals":  "readvals token used by parser/runtime flow",
}

func showSymbolInfo(app *appState) string {
	if app == nil || app.ed == nil {
		return "No symbol info"
	}
	if bufferSyntaxKind(app, app.currentPath, app.ed.Runes()) != syntaxMiranda {
		return "Symbol info: Miranda mode only"
	}

	sym := mirandaSymbolUnderCaret(app.ed.Runes(), app.ed.Caret)
	if sym == "" {
		return "No symbol under cursor"
	}
	lower := strings.ToLower(sym)

	var b strings.Builder
	fmt.Fprintf(&b, "Symbol: %s\n", sym)

	if desc, ok := mirandaKeywordInfo[lower]; ok {
		fmt.Fprintf(&b, "\nKeyword:\n%s\n", desc)
	}

	// Prefer user-local context first; this keeps Esc+I useful while editing scripts.
	if src, def, ok := mirandaDefinitionFromCurrentBuffer(app, sym); ok {
		fmt.Fprintf(&b, "\nDefinition (%s):\nCode (miranda):\n", src)
		for _, ln := range def {
			b.WriteString("    ")
			b.WriteString(strings.TrimRight(ln, "\r"))
			b.WriteRune('\n')
		}
		appendMirandaHelpNotes(&b, sym)
		return strings.TrimSpace(b.String())
	}

	// Fallback to vendored Miranda stdlib definitions for built-in/library symbols.
	if src, def, ok := mirandaDefinitionFromStdlib(sym); ok {
		fmt.Fprintf(&b, "\nDefinition (%s):\nCode (miranda):\n", src)
		for _, ln := range def {
			b.WriteString("    ")
			b.WriteString(strings.TrimRight(ln, "\r"))
			b.WriteRune('\n')
		}
		appendMirandaHelpNotes(&b, sym)
		return strings.TrimSpace(b.String())
	}

	appendMirandaHelpNotes(&b, sym)
	b.WriteString("\nNo definition found in current buffer or vendored Miranda standard library.")
	return strings.TrimSpace(b.String())
}

func mirandaDefinitionFromCurrentBuffer(app *appState, symbol string) (string, []string, bool) {
	if app == nil || app.ed == nil {
		return "", nil, false
	}
	lines := strings.Split(app.ed.String(), "\n")
	def, ok := findMirandaDefinition(lines, symbol)
	if !ok {
		return "", nil, false
	}
	src := "current buffer"
	if strings.TrimSpace(app.currentPath) != "" {
		src = app.currentPath
	}
	return src, def, true
}

func mirandaDefinitionFromStdlib(symbol string) (string, []string, bool) {
	files := []string{
		filepath.Join("miranda", "miralib", "stdenv.m"),
		filepath.Join("miranda", "miralib", "prelude"),
	}
	for _, p := range files {
		runes, err := readFileRunes(p)
		if err != nil {
			continue
		}
		def, ok := findMirandaDefinition(strings.Split(string(runes), "\n"), symbol)
		if ok {
			return p, def, true
		}
	}
	return "", nil, false
}

func findMirandaDefinition(lines []string, symbol string) ([]string, bool) {
	if strings.TrimSpace(symbol) == "" {
		return nil, false
	}
	for i := 0; i < len(lines); i++ {
		norm := normalizeMirandaCodeLine(lines[i])
		if !isMirandaSymbolDefinitionLine(norm, symbol) {
			continue
		}
		start := i
		if i > 0 {
			prev := normalizeMirandaCodeLine(lines[i-1])
			if strings.HasPrefix(prev, symbol+" ") && strings.Contains(prev, "::") {
				start = i - 1
			}
		}
		for start > 0 {
			prevRaw := strings.TrimSpace(lines[start-1])
			if strings.HasPrefix(prevRaw, "||") {
				start--
				continue
			}
			break
		}
		end := i
		for j := i + 1; j < len(lines); j++ {
			raw := lines[j]
			normJ := normalizeMirandaCodeLine(raw)
			if strings.TrimSpace(normJ) == "" {
				break
			}
			if looksLikeMirandaTopLevelDefinition(normJ) && !isMirandaSymbolDefinitionLine(normJ, symbol) {
				break
			}
			end = j
		}
		out := make([]string, 0, end-start+1)
		for k := start; k <= end; k++ {
			if strings.TrimSpace(lines[k]) == "" {
				continue
			}
			out = append(out, normalizeMirandaCodeLine(lines[k]))
		}
		if len(out) > 0 {
			return out, true
		}
	}
	return nil, false
}

func normalizeMirandaCodeLine(line string) string {
	s := strings.TrimLeft(line, " \t")
	if strings.HasPrefix(s, ">") {
		s = strings.TrimLeft(s[1:], " \t")
	}
	return s
}

func isMirandaSymbolDefinitionLine(line, symbol string) bool {
	if line == "" || symbol == "" {
		return false
	}
	if !strings.HasPrefix(line, symbol) {
		return false
	}
	if len(line) > len(symbol) {
		r, _ := utf8DecodeRuneInString(line[len(symbol):])
		if isMirandaSymbolRune(r) {
			return false
		}
	}
	return strings.Contains(line, "=") || strings.Contains(line, "::")
}

func looksLikeMirandaTopLevelDefinition(line string) bool {
	if line == "" {
		return false
	}
	r, size := utf8DecodeRuneInString(line)
	if size == 0 || !isMirandaSymbolRune(r) {
		return false
	}
	for i := size; i < len(line); {
		r2, sz := utf8DecodeRuneInString(line[i:])
		if sz == 0 {
			break
		}
		if isMirandaSymbolRune(r2) {
			i += sz
			continue
		}
		break
	}
	return strings.Contains(line, "=") || strings.Contains(line, "::")
}

func utf8DecodeRuneInString(s string) (rune, int) {
	if s == "" {
		return 0, 0
	}
	return utf8.DecodeRuneInString(s)
}

func mirandaSymbolUnderCaret(buf []rune, caret int) string {
	if caret < 0 {
		caret = 0
	}
	if caret > len(buf) {
		caret = len(buf)
	}
	if len(buf) == 0 {
		return ""
	}

	start := caret
	if start == len(buf) && start > 0 {
		start--
	}
	if start < len(buf) && !isMirandaSymbolRune(buf[start]) && start > 0 && isMirandaSymbolRune(buf[start-1]) {
		start--
	}
	if start < 0 || start >= len(buf) || !isMirandaSymbolRune(buf[start]) {
		return ""
	}

	lo := start
	for lo > 0 && isMirandaSymbolRune(buf[lo-1]) {
		lo--
	}
	hi := start
	for hi+1 < len(buf) && isMirandaSymbolRune(buf[hi+1]) {
		hi++
	}
	return string(buf[lo : hi+1])
}

func isMirandaSymbolRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '\''
}

func appendMirandaHelpNotes(b *strings.Builder, symbol string) {
	if b == nil {
		return
	}
	notes := mirandaHelpMatches(symbol, 3)
	if len(notes) == 0 {
		return
	}
	b.WriteString("\nHelp notes:\n")
	for _, h := range notes {
		b.WriteString("  ")
		b.WriteString(h)
		b.WriteRune('\n')
	}
}

func mirandaHelpMatches(symbol string, max int) []string {
	if strings.TrimSpace(symbol) == "" || max <= 0 {
		return nil
	}
	needle := strings.ToLower(symbol)
	files := []string{
		filepath.Join("miranda", "miralib", "helpfile"),
		filepath.Join("miranda", "README"),
		filepath.Join("miranda", "rules.y"),
	}
	out := make([]string, 0, max)
	seen := map[string]struct{}{}
	for _, p := range files {
		runes, err := readFileRunes(p)
		if err != nil {
			continue
		}
		lines := strings.Split(string(runes), "\n")
		for _, raw := range lines {
			line := strings.TrimSpace(raw)
			if line == "" || !strings.Contains(strings.ToLower(line), needle) {
				continue
			}
			if len(line) > 120 {
				line = line[:120] + "..."
			}
			if _, ok := seen[line]; ok {
				continue
			}
			seen[line] = struct{}{}
			out = append(out, line)
			if len(out) >= max {
				return out
			}
		}
	}
	return out
}

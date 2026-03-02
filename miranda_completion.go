package main

import (
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type mirandaCandidate struct {
	item     completionItem
	priority int
}

var (
	mirandaStdlibOnce sync.Once
	mirandaStdlibDefs map[string]completionItem
)

// mirandaIdentPrefixStart finds the symbol-start for completion at caret.
func mirandaIdentPrefixStart(buf []rune, caret int) int {
	caret = clamp(caret, 0, len(buf))
	start := caret
	for start > 0 && isMirandaSymbolRune(buf[start-1]) {
		start--
	}
	return start
}

// mirandaCompletionItems merges completion candidates from keywords, builtins,
// current buffer definitions, and vendored stdlib definitions.
func mirandaCompletionItems(app *appState, buf []rune, prefix string) []completionItem {
	if strings.TrimSpace(prefix) == "" {
		return nil
	}
	candidates := map[string]mirandaCandidate{}
	add := func(name, detail, doc string, priority int) {
		if name == "" || !strings.HasPrefix(name, prefix) {
			return
		}
		cur, ok := candidates[name]
		if ok && cur.priority >= priority {
			return
		}
		candidates[name] = mirandaCandidate{
			item: completionItem{
				Label:  name,
				Insert: name,
				Detail: detail,
				Doc:    doc,
			},
			priority: priority,
		}
	}

	for kw := range mirandaKeywords {
		add(kw, "keyword", "", 1)
	}
	for bi := range mirandaBuiltins {
		add(bi, "builtin", "", 2)
	}
	localDefs := mirandaLocalDefinitions(app, buf)
	for name, sig := range localDefs {
		add(name, "current file", sig, 4)
	}
	for name, item := range mirandaStdlibCompletionItems() {
		add(name, item.Detail, item.Doc, 3)
	}

	names := make([]string, 0, len(candidates))
	for name := range candidates {
		names = append(names, name)
	}
	sort.Strings(names)

	items := make([]completionItem, 0, len(names))
	for _, name := range names {
		items = append(items, candidates[name].item)
	}
	return items
}

func mirandaLocalDefinitions(app *appState, buf []rune) map[string]string {
	if app == nil || app.bufIdx < 0 || app.bufIdx >= len(app.buffers) {
		return mirandaDefinitionsFromLines(strings.Split(string(buf), "\n"))
	}
	slot := &app.buffers[app.bufIdx]
	if slot.completionDefs != nil && slot.completionDefsTextRev == slot.textRev {
		return slot.completionDefs
	}
	defs := mirandaDefinitionsFromLines(strings.Split(string(buf), "\n"))
	slot.completionDefs = defs
	slot.completionDefsTextRev = slot.textRev
	return defs
}

// mirandaStdlibCompletionItems indexes stdlib names once per process.
func mirandaStdlibCompletionItems() map[string]completionItem {
	mirandaStdlibOnce.Do(func() {
		mirandaStdlibDefs = map[string]completionItem{}
		base, ok := findVendoredMirandaDir()
		if !ok {
			return
		}
		files := []string{
			filepath.Join(base, "miralib", "stdenv.m"),
			filepath.Join(base, "miralib", "prelude"),
		}
		for _, p := range files {
			runes, err := readFileRunes(p)
			if err != nil {
				continue
			}
			for name, sig := range mirandaDefinitionsFromLines(strings.Split(string(runes), "\n")) {
				if _, exists := mirandaStdlibDefs[name]; exists {
					continue
				}
				mirandaStdlibDefs[name] = completionItem{
					Label:  name,
					Insert: name,
					Detail: "stdlib (" + filepath.Base(p) + ")",
					Doc:    sig,
				}
			}
		}
	})
	return mirandaStdlibDefs
}

// mirandaDefinitionsFromLines collects top-level definition names/signatures.
func mirandaDefinitionsFromLines(lines []string) map[string]string {
	out := map[string]string{}
	for _, line := range lines {
		norm := normalizeMirandaCodeLine(line)
		if !looksLikeMirandaTopLevelDefinition(norm) {
			continue
		}
		name, ok := mirandaDefinitionName(norm)
		if !ok || name == "" {
			continue
		}
		if _, exists := out[name]; exists {
			continue
		}
		out[name] = norm
	}
	return out
}

func mirandaDefinitionName(line string) (string, bool) {
	if line == "" {
		return "", false
	}
	r, size := utf8DecodeRuneInString(line)
	if size == 0 || !isMirandaSymbolRune(r) {
		return "", false
	}
	i := size
	for i < len(line) {
		next, n := utf8DecodeRuneInString(line[i:])
		if n == 0 || !isMirandaSymbolRune(next) {
			break
		}
		i += n
	}
	name := line[:i]
	if name == "" {
		return "", false
	}
	return name, true
}

// longestSharedPrefix returns the maximal shared insert prefix across items.
func longestSharedPrefix(items []completionItem) string {
	if len(items) == 0 {
		return ""
	}
	base := items[0].Insert
	if base == "" {
		base = items[0].Label
	}
	for i := 1; i < len(items); i++ {
		cur := items[i].Insert
		if cur == "" {
			cur = items[i].Label
		}
		base = sharedPrefix(base, cur)
		if base == "" {
			return ""
		}
	}
	return base
}

func sharedPrefix(a, b string) string {
	ar := []rune(a)
	br := []rune(b)
	n := min(len(ar), len(br))
	i := 0
	for i < n && ar[i] == br[i] {
		i++
	}
	return string(ar[:i])
}

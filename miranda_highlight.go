package main

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

var mirandaKeywords = map[string]struct{}{
	"abstype": {}, "div": {}, "if": {}, "include": {}, "otherwise": {},
	"readvals": {}, "show": {}, "type": {}, "where": {}, "with": {},
	"mod": {}, "free": {}, "export": {}, "bnf": {}, "lex": {},
}

var mirandaBuiltins = map[string]struct{}{
	"hd": {}, "tl": {}, "map": {}, "filter": {}, "foldl": {}, "foldr": {},
	"zip": {}, "zip2": {}, "take": {}, "drop": {}, "member": {}, "reverse": {},
}

func (h *syntaxHighlighter) lineStyleForKind(path string, textRev int, lines []string, kind syntaxKind) [][]tokenStyle {
	if h == nil || len(lines) == 0 || kind == syntaxNone {
		return nil
	}
	if h.lastPath == path && h.lastTextRev == textRev && h.lastKind == kind {
		return h.lineStyles
	}

	var out [][]tokenStyle
	switch kind {
	case syntaxMarkdown:
		out = highlightMarkdownLines(lines)
	case syntaxMiranda:
		out = highlightMirandaLines(lines)
	default:
		out = nil
	}

	h.lastPath = path
	h.lastTextRev = textRev
	h.lastKind = kind
	h.lineStyles = out
	return out
}

func highlightMarkdownLines(lines []string) [][]tokenStyle {
	out := make([][]tokenStyle, len(lines))
	hasAny := false
	for i, line := range lines {
		rs := []rune(line)
		if len(rs) == 0 {
			continue
		}
		styles := make([]tokenStyle, len(rs))
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "#") {
			for j := range styles {
				styles[j] = styleHeading
			}
			hasAny = true
			out[i] = styles
			continue
		}
		for j := 0; j < len(rs); j++ {
			if rs[j] == '`' {
				k := j + 1
				for k < len(rs) && rs[k] != '`' {
					styles[k] = styleString
					k++
				}
				hasAny = true
				j = k
				continue
			}
			if rs[j] == '[' {
				k := j
				for k < len(rs) && rs[k] != ')' {
					styles[k] = styleLink
					k++
				}
				if k < len(rs) {
					styles[k] = styleLink
				}
				hasAny = true
				j = k
			}
		}
		if rowHasStyle(styles) {
			out[i] = styles
		}
	}
	if !hasAny {
		return nil
	}
	return out
}

func highlightMirandaLines(lines []string) [][]tokenStyle {
	out := make([][]tokenStyle, len(lines))
	hasAny := false
	for i, line := range lines {
		rs := []rune(line)
		if len(rs) == 0 {
			continue
		}
		styles := make([]tokenStyle, len(rs))

		for j := 0; j < len(rs); {
			if j+1 < len(rs) && rs[j] == '|' && rs[j+1] == '|' {
				for k := j; k < len(rs); k++ {
					styles[k] = styleComment
				}
				hasAny = true
				break
			}
			if rs[j] == '"' {
				styles[j] = styleString
				k := j + 1
				for k < len(rs) {
					styles[k] = styleString
					if rs[k] == '"' && rs[k-1] != '\\' {
						k++
						break
					}
					k++
				}
				hasAny = true
				j = k
				continue
			}
			if rs[j] == '\'' {
				styles[j] = styleString
				k := j + 1
				for k < len(rs) {
					styles[k] = styleString
					if rs[k] == '\'' && rs[k-1] != '\\' {
						k++
						break
					}
					k++
				}
				hasAny = true
				j = k
				continue
			}
			if unicode.IsDigit(rs[j]) {
				k := j + 1
				for k < len(rs) && (unicode.IsDigit(rs[k]) || rs[k] == '.') {
					k++
				}
				for x := j; x < k; x++ {
					styles[x] = styleNumber
				}
				hasAny = true
				j = k
				continue
			}
			if isMirandaIdentStart(rs[j]) {
				k := j + 1
				for k < len(rs) && isMirandaIdentRune(rs[k]) {
					k++
				}
				ident := string(rs[j:k])
				identLower := strings.ToLower(ident)
				style := styleDefault
				if _, ok := mirandaKeywords[identLower]; ok {
					style = styleKeyword
				} else if _, ok := mirandaBuiltins[identLower]; ok {
					style = styleFunction
				} else if unicode.IsUpper(rs[j]) {
					style = styleType
				}
				if style != styleDefault {
					for x := j; x < k; x++ {
						styles[x] = style
					}
					hasAny = true
				}
				j = k
				continue
			}
			if isMirandaPunctuation(rs[j]) {
				styles[j] = stylePunctuation
				hasAny = true
			}
			j++
		}
		if rowHasStyle(styles) {
			out[i] = styles
		}
	}
	if !hasAny {
		return nil
	}
	return out
}

func isMirandaIdentStart(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isMirandaIdentRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '\''
}

func isMirandaPunctuation(r rune) bool {
	return strings.ContainsRune("=:+-*/<>|&!.,()[]{}\\", r)
}

func rowHasStyle(row []tokenStyle) bool {
	for _, s := range row {
		if s != styleDefault {
			return true
		}
	}
	return false
}

func runeColFromByte(s string, b int) int {
	if b <= 0 {
		return 0
	}
	if b >= len(s) {
		return utf8.RuneCountInString(s)
	}
	return utf8.RuneCountInString(s[:b])
}

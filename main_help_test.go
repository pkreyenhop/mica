package main

import (
	"strings"
	"testing"
)

func TestHelpTextIncludesAllEntries(t *testing.T) {
	h := helpText()
	if strings.TrimSpace(h) == "" {
		t.Fatal("help text should not be empty")
	}
	for _, e := range helpEntries {
		if !strings.Contains(h, e.action) || !strings.Contains(h, e.keys) {
			t.Fatalf("help text missing %q / %q", e.action, e.keys)
		}
	}
}

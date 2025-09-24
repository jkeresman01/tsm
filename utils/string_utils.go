package utils

import (
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"
)

var MatchStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

func HighlightMatches(item, query string) string {
	if query == "" {
		return item
	}
	qset := RuneSetFold(query)
	var b strings.Builder
	for _, r := range item {
		if _, ok := qset[unicode.ToLower(r)]; ok {
			b.WriteString(MatchStyle.Render(string(r)))
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func FuzzyFilter(items []string, query string) []string {
	if query == "" {
		return items
	}
	q := strings.ToLower(query)
	out := make([]string, 0, len(items))
	for _, it := range items {
		if strings.Contains(strings.ToLower(it), q) {
			out = append(out, it)
		}
	}
	return out
}

func RuneSetFold(s string) map[rune]struct{} {
	m := make(map[rune]struct{}, len(s))
	for _, r := range s {
		m[unicode.ToLower(r)] = struct{}{}
	}
	return m
}

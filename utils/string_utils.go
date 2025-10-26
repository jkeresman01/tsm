package utils

import (
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"
)

// ///////////////////////////////////////////////////////////////////////////////////////////
//
// MatchStyle defines the style for highlighting matched characters.
//
// ///////////////////////////////////////////////////////////////////////////////////////////
var MatchStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			HighlightMatches highlights characters in item that match the query.
//
//		@Description	Uses MatchStyle to highlight matching characters
//
//		@Param			item	string	String to highlight
//		@Param			query	string	Search query for matching
//
//		@Return			string	String with highlighted matches
//
// ///////////////////////////////////////////////////////////////////////////////////////////
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

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			FuzzyFilter filters items based on fuzzy string matching.
//
//		@Description	Returns items where query appears as substring (case-insensitive)
//
//		@Param			items	[]string	Items to filter
//		@Param			query	string		Search query
//
//		@Return			[]string	Filtered items
//
// ///////////////////////////////////////////////////////////////////////////////////////////
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

// ///////////////////////////////////////////////////////////////////////////////////////////
//
//	 @Brief			RuneSetFold creates a set of lowercase runes from a string.
//
//		@Description	Used for case-insensitive character matching
//
//		@Param			s	string	Input string
//
//		@Return			map[rune]struct{}	Set of lowercase runes
//
// ///////////////////////////////////////////////////////////////////////////////////////////
func RuneSetFold(s string) map[rune]struct{} {
	m := make(map[rune]struct{}, len(s))
	for _, r := range s {
		m[unicode.ToLower(r)] = struct{}{}
	}
	return m
}

//go:generate go run generate.go

// Package coname normalizes company names by stripping legal entity
// designators (Ltd, Inc, GmbH, S.A., etc.) and common prefixes like "The".
// It supports 500+ designators from 20+ countries including English, German,
// French, Spanish, Italian, Dutch, Nordic, Polish, Russian, Japanese, Korean,
// Chinese, and more.
//
// Use [NormalizeWords] to clean company names for matching, deduplication,
// or display.
package coname

import (
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

// prefixes to strip (e.g. "The Coca-Cola Company" → "Coca-Cola Company" → "Coca-Cola")
var prefixes = []string{
	"the",
}

// trimCutset is the punctuation and whitespace trimmed where a designator
// was removed.
const trimCutset = " \t\n\v\f\r,.-&'/\\"

// Compiled designator lists: normalized, deduplicated, and sorted
// longest-first so greedy matching prefers "sendirian berhad" over "berhad".
var (
	normSuffixes = compile(generatedSuffixes)
	normPrefixes = compile(prefixes)
)

func compile(designators []string) []string {
	seen := make(map[string]bool, len(designators))
	out := make([]string, 0, len(designators))
	for _, d := range designators {
		d = strings.ReplaceAll(strings.ToLower(d), ".", "")
		if !seen[d] {
			seen[d] = true
			out = append(out, d)
		}
	}
	slices.SortStableFunc(out, func(a, b string) int { return len(b) - len(a) })
	return out
}

// matchSuffix reports where the compiled designator d (lowercase, no periods)
// matches the end of name, comparing case-insensitively and skipping '.'
// runes in name. Returns the byte offset in name where the match starts,
// or -1 if it doesn't match.
func matchSuffix(name, d string) int {
	i, j := len(name), len(d)
	for j > 0 {
		if i == 0 {
			return -1
		}
		r, size := utf8.DecodeLastRuneInString(name[:i])
		if r == '.' {
			i -= size
			continue
		}
		dr, dsize := utf8.DecodeLastRuneInString(d[:j])
		if unicode.ToLower(r) != dr {
			return -1
		}
		i -= size
		j -= dsize
	}
	return i
}

// matchPrefix reports where the compiled designator d matches the start of
// name (case-insensitive, '.'-skipping), followed by a space and at least one
// more non-period rune. Returns the byte offset just past the space, or -1.
func matchPrefix(name, d string) int {
	i, j := 0, 0
	for j < len(d) {
		if i >= len(name) {
			return -1
		}
		r, size := utf8.DecodeRuneInString(name[i:])
		if r == '.' {
			i += size
			continue
		}
		dr, dsize := utf8.DecodeRuneInString(d[j:])
		if unicode.ToLower(r) != dr {
			return -1
		}
		i += size
		j += dsize
	}
	for {
		if i >= len(name) {
			return -1
		}
		r, size := utf8.DecodeRuneInString(name[i:])
		if r == '.' {
			i += size
			continue
		}
		if r != ' ' {
			return -1
		}
		i += size
		break
	}
	rest := name[i:]
	for {
		if rest == "" {
			return -1
		}
		r, size := utf8.DecodeRuneInString(rest)
		if r != '.' {
			break
		}
		rest = rest[size:]
	}
	return i
}

// lastNonPeriod returns the last rune of s that isn't '.', or 0 if none.
func lastNonPeriod(s string) rune {
	for len(s) > 0 {
		r, size := utf8.DecodeLastRuneInString(s)
		if r != '.' {
			return r
		}
		s = s[:len(s)-size]
	}
	return 0
}

// normalize strips legal entity designators from a company name and returns
// the cleaned name. Case-insensitive and period-insensitive. Never strips
// the entire name. The result is always a substring of the input; normalize
// performs no allocations.
func normalize(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}

	// strip suffix designators. Loop because some names stack several,
	// e.g. "Foo GmbH & Co. KG".
	for stripped := true; stripped; {
		stripped = false
		for _, d := range normSuffixes {
			cut := matchSuffix(name, d)
			if cut < 0 {
				continue
			}
			// skips strip the whole name, and require a word boundary
			// so "Tabasco" doesn't lose its "co". lastNonPeriod skips
			// dots so "X.LLC" keeps the same boundary semantics as the
			// old normalized-buffer version.
			if r := lastNonPeriod(name[:cut]); r == 0 || isWordChar(r) {
				continue
			}
			name = strings.TrimRight(name[:cut], trimCutset)
			stripped = true
			break
		}
	}

	for _, d := range normPrefixes {
		if i := matchPrefix(name, d); i >= 0 {
			name = strings.TrimLeft(name[i:], trimCutset)
		}
	}

	return strings.TrimSpace(name)
}

// NormalizeWords strips legal entity designators from a company name and
// collapses each run of whitespace to a single space. Matching is
// case-insensitive and period-insensitive (so "L.L.C." and "LLC" are treated
// the same). It will never strip the entire name, protecting against short
// names that collide with designator abbreviations.
//
//	NormalizeWords("Acme Ltd.")        → "Acme"
//	NormalizeWords("THE COCA-COLA CO") → "COCA-COLA"
//	NormalizeWords("Coca  Cola Co")    → "Coca Cola"
//	NormalizeWords("SAP SE")           → "SAP"
//
// Unless whitespace needs collapsing, the result is a substring of the input
// and no copy is made.
func NormalizeWords(name string) string {
	n := normalize(name)
	if singleSpaced(n) {
		return n
	}
	return strings.Join(strings.Fields(n), " ")
}

// singleSpaced reports whether the only whitespace in s is single ASCII spaces.
func singleSpaced(s string) bool {
	prevSpace := false
	for _, r := range s {
		if !unicode.IsSpace(r) {
			prevSpace = false
			continue
		}
		if r != ' ' || prevSpace {
			return false
		}
		prevSpace = true
	}
	return true
}

func isWordChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

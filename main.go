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
	"bytes"
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

// appendNormalized appends the lowercased form of s, minus periods, to norm.
// For every byte appended it records in src the byte offset in s of the rune
// it came from, so match positions in the normalized form can be mapped back
// to s even when their lengths differ (dropped periods, case-changed runes).
func appendNormalized(norm []byte, src []int, s string) ([]byte, []int) {
	for i, r := range s {
		if r == '.' {
			continue
		}
		n := len(norm)
		norm = utf8.AppendRune(norm, unicode.ToLower(r))
		for ; n < len(norm); n++ {
			src = append(src, i)
		}
	}
	return norm, src
}

// normalize strips legal entity designators from a company name and returns
// the cleaned name. It is case-insensitive and period-insensitive (so
// "L.L.C." and "LLC" are treated the same). It will never strip the entire
// name, protecting against short names that collide with designator
// abbreviations (e.g., "SAP" won't be stripped even though it matches "S.A.P.").
//
// The result is a substring of the input; normalize never copies the name.
func normalize(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}

	norm, src := appendNormalized(make([]byte, 0, len(name)), make([]int, 0, len(name)), name)

	// Strip suffix designators. Loop because some names stack several,
	// e.g. "Foo GmbH & Co. KG".
	for stripped := true; stripped; {
		stripped = false
		for _, d := range normSuffixes {
			cut := len(norm) - len(d)
			if cut < 0 || string(norm[cut:]) != d {
				continue
			}
			if cut == 0 {
				continue // don't strip if it would remove everything
			}
			// Require a word boundary so "Tabasco" doesn't lose its "co".
			if cut > 0 {
				if r, _ := utf8.DecodeLastRune(norm[:cut]); isWordChar(r) {
					continue
				}
			}
			name = strings.TrimRight(name[:src[cut]], trimCutset)
			norm = bytes.TrimRight(norm[:cut], trimCutset)
			stripped = true
			break
		}
	}

	for _, d := range normPrefixes {
		if len(norm) > len(d)+1 && string(norm[:len(d)]) == d && norm[len(d)] == ' ' {
			name = strings.TrimLeft(name[src[len(d)]:], trimCutset)
			norm, src = appendNormalized(norm[:0], src[:0], name)
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

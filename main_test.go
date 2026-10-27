package main

import (
	"slices"
	"testing"
)

// Each case targets one rule from the requirements. The comment above each
// group names the kind of edge case, so you can reuse the categories later.
func TestProcessInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		// Baseline: the simplest input that should obviously work.
		{"single word", "hello", []string{"hello"}},
		{"multiple words", "the cat sat", []string{"the", "cat", "sat"}},

		// Empty and "nothing useful" input: boundaries of size zero.
		{"empty input", "", nil},
		{"only whitespace", "   \n\t  ", nil},
		{"only punctuation", "!!! ... ---", nil},

		// Normalization: different spellings that must become the same word.
		{"mixed case", "Hello HELLO hello", []string{"hello", "hello", "hello"}},
		{"curly apostrophe matches straight", "don’t don't", []string{"don't", "don't"}},

		// Separators: everything that isn't part of a word must split words.
		{"punctuation splits", "end. start,stop", []string{"end", "start", "stop"}},
		{"hyphen splits", "post-game", []string{"post", "game"}},
		{"newlines and tabs split", "one\ntwo\tthree", []string{"one", "two", "three"}},

		// Apostrophes: kept only between letters/digits (decision b).
		{"contraction kept", "That's", []string{"that's"}},
		{"multiple internal apostrophes", "rock'n'roll", []string{"rock'n'roll"}},
		{"quoting apostrophes stripped", "'hello'", []string{"hello"}},
		{"curly quotes stripped", "‘quoted’", []string{"quoted"}},
		{"trailing possessive stripped", "dogs'", []string{"dogs"}},

		// Character classes: digits and non-ASCII letters are part of words.
		{"digits are words", "2026", []string{"2026"}},
		{"letters and digits mixed", "v2 1990s", []string{"v2", "1990s"}},
		{"non-ASCII letters", "café naïve", []string{"café", "naïve"}},
		{"non-ASCII uppercase lowered", "ÉCOLE", []string{"école"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProcessInput([]byte(tt.input))
			if !slices.Equal(got, tt.want) {
				t.Errorf("ProcessInput(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestCountWordFreq(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []WordCount
	}{
		// Size boundaries: zero and one.
		{"no words", nil, nil},
		{"one word", []string{"a"}, []WordCount{{"a", 1}}},

		// Ordering rules: higher count first, then alphabetical for ties.
		{"sorted by count descending", []string{"b", "a", "b"}, []WordCount{{"b", 2}, {"a", 1}}},
		{"ties broken alphabetically", []string{"c", "a", "b"}, []WordCount{{"a", 1}, {"b", 1}, {"c", 1}}},
		{"count beats alphabet", []string{"z", "z", "a"}, []WordCount{{"z", 2}, {"a", 1}}},

		// Mixed: several counts with ties inside each count group.
		{"ties within count groups", []string{"d", "b", "d", "a", "c", "c"},
			[]WordCount{{"c", 2}, {"d", 2}, {"a", 1}, {"b", 1}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CountWordFreq(tt.input)
			if !slices.Equal(got, tt.want) {
				t.Errorf("CountWordFreq(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestTopN(t *testing.T) {
	three := []WordCount{{"a", 3}, {"b", 2}, {"c", 1}}

	tests := []struct {
		name  string
		input []WordCount
		n     int
		want  []WordCount
	}{
		// Boundaries around len(input): less than, equal to, greater than.
		{"n smaller than length", three, 2, []WordCount{{"a", 3}, {"b", 2}}},
		{"n equal to length", three, 3, three},
		{"n larger than length", three, 10, three},

		// Smallest valid n, and empty input.
		{"n is one", three, 1, []WordCount{{"a", 3}}},
		{"empty input", nil, 5, nil},

		// Invalid n: n <= 0 means "no limit" and returns everything.
		{"n is zero", three, 0, three},
		{"n is negative", three, -1, three},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TopN(tt.input, tt.n)
			if !slices.Equal(got, tt.want) {
				t.Errorf("TopN(%v, %d) = %v, want %v", tt.input, tt.n, got, tt.want)
			}
		})
	}
}

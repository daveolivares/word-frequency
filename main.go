package main

import (
	"cmp"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"strings"
)

var wordRegex = regexp.MustCompile(`[\p{L}\p{N}]+(?:'[\p{L}\p{N}]+)*`)
var replacer = strings.NewReplacer(
	"’", "'", // Right single quote (Most common smart apostrophe)
	"‘", "'", // Left single quote
	"＇", "'", // Full-width apostrophe
)

type WordCount struct {
	Word  string
	Count int
}

func ProcessInput(data []byte) []string {
	// Convert to lower case
	text := strings.ToLower(string(data))

	text = replacer.Replace(text)

	words := wordRegex.FindAllString(text, -1)

	return words
}

func CountWordFreq(words []string) []WordCount {
	// Count word frequency
	wordCounts := make(map[string]int)
	for _, word := range words {
		wordCounts[word]++
	}

	// Build the result slice
	result := make([]WordCount, 0, len(wordCounts))
	for word, count := range wordCounts {
		result = append(result, WordCount{Word: word, Count: count})
	}

	// Sort based on WordCount.Count first
	// then WordCount.Word for ties
	slices.SortFunc(result, func(a, b WordCount) int {
		if n := cmp.Compare(b.Count, a.Count); n != 0 {
			return n
		}
		return cmp.Compare(a.Word, b.Word)
	})

	return result
}

// TopN returns the first n entries of wordCounts. If n <= 0 or n exceeds
// the length, it returns all entries.
func TopN(wordCounts []WordCount, nFlag int) []WordCount {
	if nFlag <= 0 || nFlag > len(wordCounts) {
		nFlag = len(wordCounts)
	}

	return wordCounts[:nFlag]
}

func PrintWordCounts(wordCounts []WordCount) {
	for _, wc := range wordCounts {
		fmt.Printf("%s\t%d\n", wc.Word, wc.Count)
	}
}

func main() {
	// 1. Define flag(s)
	nFlag := flag.Int("n", 10, "print the first N results")

	// 2. Parse flags first
	flag.Parse()

	// 3. Validate nFlag
	if *nFlag <= 0 {
		fmt.Fprintln(os.Stderr, "error: -n must be > 0")
		os.Exit(1)
	}

	// 4. Read data from stdin
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
		os.Exit(1)
	}

	// 5. Process data
	words := ProcessInput(data)
	result := CountWordFreq(words)
	PrintWordCounts(TopN(result, *nFlag))
}

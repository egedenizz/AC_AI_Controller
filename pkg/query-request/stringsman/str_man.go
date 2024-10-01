package stringsman

import (
	"regexp"
	"strings"
)

func FindNumberInString(input string) (string, string, string) {

	re := regexp.MustCompile(`\d+`)

	loc := re.FindStringIndex(input)

	if loc == nil {
		return "", "", ""
	}

	before := input[:loc[0]]
	number := input[loc[0]:loc[1]]
	after := input[loc[1]:]

	return before, number, after
}

func SplitByWords(input string, words []string) []string {

	wordMap := make(map[string]struct{})
	for _, word := range words {
		wordMap[word] = struct{}{}
	}

	var parts []string

	start := 0
	for i := 0; i < len(input); i++ {
		for _, word := range words {
			if strings.HasPrefix(strings.ToLower(input[i:]), word) {
				if i > start {
					parts = append(parts, input[start:i])
				}
				parts = append(parts, word)
				start = i + len(word)
				i = start - 1
				break
			}
		}
	}

	if start < len(input) {
		parts = append(parts, input[start:])
	}

	return parts
}

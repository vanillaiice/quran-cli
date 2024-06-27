package arabic

import "strings"

// vowels is a list of arabic vowels.
var vowels = []rune{
	'\u064E', // Fatha
	'\u0650', // Kasra
	'\u064F', // Damma
	'\u0627', // Alif
	'\u064A', // Ya
	'\u0648', // Waw
	'\u0652', // Sukun
	'\u0651', // Shadda
}

// WrapToSlice wraps a string according to the limit
// and returns it to a slice of strings.
func WrapToSlice(s string, limit int) []string {
	s = Wrap(s, limit)
	return strings.Split(s, "\n")
}

// Wrap wraps a string according to the limit.
func Wrap(s string, limit int) string {
	wrapped := []rune{}

	var i int
	for _, c := range s {
		if i == limit {
			i = 0
			wrapped = append(wrapped, '\n')
		}

		if IsVowel(c) {
			wrapped = append(wrapped, c)
			continue
		}

		wrapped = append(wrapped, c)
		i++
	}

	return string(wrapped)
}

// Count returns the number of arabic letters in a string.
func Count(s string) int {
	var count int

	for _, c := range s {
		if IsVowel(c) {
			continue
		}
		count++
	}

	return count
}

// IsVowel returns true if a rune is an arabic vowel.
func IsVowel(r rune) bool {
	for _, v := range vowels {
		if v == r {
			return true
		}
	}
	return false
}

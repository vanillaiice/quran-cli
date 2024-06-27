package arabic

import (
	"strconv"
)

// numMap is a map that maps a number to its arabic script notation.
var numMap map[rune]rune = map[rune]rune{
	'0': '٠',
	'1': '١',
	'2': '٢',
	'3': '٣',
	'4': '٤',
	'5': '٥',
	'6': '٦',
	'7': '٧',
	'8': '٨',
	'9': '٩',
}

// ToArabic converts a number to its arabic script notation.
func ToArabic(num int) (n string) {
	numStr := strconv.Itoa(num)

	numRune := make([]rune, len(numStr))
	for i, s := range numStr {
		numRune[i] = numMap[s]
	}

	return string(numRune)
}

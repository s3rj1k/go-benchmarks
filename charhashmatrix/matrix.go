package charhashmatrix

import (
	"unicode"

	xxhash "github.com/cespare/xxhash/v2"
)

/*
	HashMatrix [y][x]bool
		y - signifies the positions within strings
		x - used for the character set
*/

type HashMatrix [][]bool

var (
	// lettersCount represents the total number of lowercase alphabetic characters ('a' to 'z')
	lettersCount = int('z'-'a') + 1
	// digitsCount represents the total number of numeric digits ('0' to '9')
	digitsCount = int('9'-'0') + 1

	totalCharactersCount = lettersCount + digitsCount
)

func uint64ToHexRunes(val uint64) []rune {
	if val == 0 {
		return []rune{'0'}
	}

	// Calculate the length of the hexadecimal representation.
	length := 0
	for v := val; v > 0; v /= 16 {
		length++
	}

	// Allocate a slice of runes with the exact required length.
	runes := make([]rune, length)

	// Fill the slice from the end.
	for i := length - 1; i >= 0; i-- {
		digit := val % 16

		if digit < 10 {
			runes[i] = rune('0' + digit)
		} else {
			runes[i] = rune('a' + (digit - 10))
		}

		val /= 16
	}

	return runes
}

func charToIndex(c rune) int {
	c = unicode.ToLower(c)

	switch {
	case '0' <= c && c <= '9':
		return int(c - '0')
	case 'a' <= c && c <= 'z':
		return digitsCount + int(c-'a')
	default:
		panic("invalid rune for char to index conversation")
	}
}

func indexToChar(i int) rune {
	switch {
	case i < digitsCount:
		// Calculate the character for digits '0'-'9'.
		return rune('0' + i)
	case i >= digitsCount:
		// Calculate the character for letters 'a'-'z'.
		return rune('a' + (i - digitsCount))
	default:
		panic("invalid index for index to char conversation")
	}
}

func hasher(s string) []rune {
	return uint64ToHexRunes(xxhash.Sum64String(s))
}

func NewMatrix() *HashMatrix {
	matrix := make(HashMatrix, 16)

	for y := range matrix {
		matrix[y] = make([]bool, y+1)

		for x := range matrix[y] {
			matrix[x] = make([]bool, totalCharactersCount)
		}
	}

	return &matrix
}

func (m *HashMatrix) hasSetCount(y int) int {
	count := 0

	for _, isSet := range (*m)[y] {
		if isSet {
			count++
		}
	}

	return count
}

func (m *HashMatrix) Set(s string) error {
	if len(s) == 0 {
		return nil
	}

	runes := hasher(s)

	for y, char := range runes {
		x := charToIndex(char)

		(*m)[y][x] = true
	}

	return nil
}

func (m *HashMatrix) Contains(s string) bool {
	if len(s) == 0 {
		return false
	}

	runes := hasher(s)

	// Loop over all rows from the end to the start.
	for y := len(runes) - 1; y >= 0; y-- {
		char := runes[y]

		x := charToIndex(char)
		if !(*m)[y][x] {
			return false
		}
	}

	return true
}

func (m *HashMatrix) Unset(s string) error {
	if len(s) == 0 {
		return nil
	}

	runes := hasher(s)

	y := len(runes) - 1

	// Initially unset the last character in the input slice.
	x := charToIndex(runes[y])

	(*m)[y][x] = false

	// Loop over all rows from the end to the start, including the last one, excluding first.
	for y := len(runes) - 1; y > 0; y-- {
		// Check if the current row has a set count of 0.
		if m.hasSetCount(y) > 0 {
			continue
		}

		// Convert the character in the row above to its index, if current row is empty.
		x := charToIndex(runes[y-1])

		// Unset the character in the row above if the current row has a set count of 0.
		(*m)[y-1][x] = false
	}

	return nil
}

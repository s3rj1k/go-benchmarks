package charbyteshashmatrix

import (
	"unicode"

	xxhash "github.com/cespare/xxhash/v2"
)

/*
	HashMatrix [y][x]bool
		y - signifies the positions within strings
		x - used for the character set
*/

type HashMatrix [][]byte

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
		matrix[y] = make([]byte, (totalCharactersCount+7)/8) // Allocate enough bytes to cover all characters
	}

	return &matrix
}

func (m *HashMatrix) setBit(y int, x int) {
	(*m)[y][x/8] |= 1 << (x % 8)
}

func (m *HashMatrix) clearBit(y int, x int) {
	(*m)[y][x/8] &^= 1 << (x % 8)
}

func (m *HashMatrix) isBitSet(y int, x int) bool {
	return ((*m)[y][x/8] & (1 << (x % 8))) != 0
}

func (*HashMatrix) ifZero(bytes []byte) bool {
	for _, v := range bytes {
		if v != 0 {
			return false
		}
	}

	return true
}

func (m *HashMatrix) Set(s string) error {
	if len(s) == 0 {
		return nil
	}

	runes := hasher(s)

	for y, char := range runes {
		x := charToIndex(char)
		m.setBit(y, x)
	}

	return nil
}

func (m *HashMatrix) Contains(s string) bool {
	if len(s) == 0 {
		return false
	}

	runes := hasher(s)

	for y := len(runes) - 1; y >= 0; y-- {
		char := runes[y]

		x := charToIndex(char)
		if !m.isBitSet(y, x) {
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

	// Unset the bit for the last character in the hashed string.
	if len(runes) > 0 {
		lastCharIndex := len(runes) - 1
		x := charToIndex(runes[lastCharIndex])
		m.clearBit(lastCharIndex, x)
	}

	// Loop over all rows from the end to the start.
	for y := len(runes) - 1; y >= 0; y-- {
		if m.ifZero((*m)[y]) && y > 0 {
			previousCharIndex := charToIndex(runes[y-1])
			m.clearBit(y-1, previousCharIndex)
		}
	}

	return nil
}

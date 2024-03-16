package charmatrix3d

import (
	"errors"
	"fmt"
	"unicode"
)

/*
	CharMatrix [z][y][x]bool
		z - represents the depth, corresponding to different string lengths
		y - signifies the positions within strings of a given length z
		x - used for the character set
*/

type CharMatrix [][][]bool

var (
	// lettersCount represents the total number of lowercase alphabetic characters ('a' to 'z')
	lettersCount = int('z'-'a') + 1
	// digitsCount represents the total number of numeric digits ('0' to '9')
	digitsCount = int('9'-'0') + 1
	// symbolsCount represents the total number of specific symbol characters ('-', '.', '/')
	symbolsCount = 3

	totalCharactersCount = lettersCount + digitsCount + symbolsCount
)

var (
	ErrInvalidCharacter = errors.New("invalid character")
	ErrIndexOutOfRange  = errors.New("index out of range")
	ErrInvalidLength    = errors.New("invalid length")
)

func charToIndex(c rune) (int, error) {
	c = unicode.ToLower(c)

	switch {
	case c == '-':
		return 0, nil
	case c == '.':
		return 1, nil
	case c == '/':
		return 2, nil
	case '0' <= c && c <= '9':
		return symbolsCount + int(c-'0'), nil
	case 'a' <= c && c <= 'z':
		return symbolsCount + digitsCount + int(c-'a'), nil
	default:
		return -1, ErrInvalidCharacter
	}
}

func indexToChar(i int) (rune, error) {
	switch {
	case i == 0:
		return '-', nil
	case i == 1:
		return '.', nil
	case i == 2:
		return '/', nil
	case i >= symbolsCount && i < symbolsCount+digitsCount:
		// Calculate the character for digits '0'-'9'.
		return rune('0' + (i - symbolsCount)), nil
	case i >= symbolsCount+digitsCount:
		// Calculate the character for letters 'a'-'z'.
		return rune('a' + (i - symbolsCount - digitsCount)), nil
	default:
		return rune(-1), ErrIndexOutOfRange
	}
}

func NewMatrix(maxStrLen int) *CharMatrix {
	// Create a 3D matrix with 'maxStrLen' layers, where each layer corresponds
	// to strings of different lengths (from 1 to maxStrLen).
	matrix := make(CharMatrix, maxStrLen)

	for z := range matrix {
		// For each string length 'z+1', initialize a 2D slice. The '+1' accounts
		// for the fact that 'z' is zero-based, while string lengths are 1-based.
		matrix[z] = make([][]bool, z+1)

		for y := range matrix[z] {
			// For each position 'y' in a string of length 'z+1', initialize a slice
			// to represent the presence or absence of each character in the character set.
			matrix[z][y] = make([]bool, totalCharactersCount)
		}
	}

	return &matrix
}

func (m *CharMatrix) hasSetCount(z, y int) int {
	count := 0

	for _, isSet := range (*m)[z][y] {
		if isSet {
			count++
		}
	}

	return count
}

func (m *CharMatrix) calculateMaxIndexLength(depth int) int {
	var maxRowIndexLength int

	// Calculate the length of the largest character set index
	maxIndexLength := len(fmt.Sprintf("%d", totalCharactersCount-1))

	if depth >= 0 && depth < len((*m)) {
		// Calculate the length of the largest row index at the given depth
		maxRowIndexLength = len(fmt.Sprintf("%d", len((*m)[depth])-1))
	}

	if maxRowIndexLength > maxIndexLength {
		return maxRowIndexLength
	}

	return maxIndexLength
}

func (m *CharMatrix) Set(s []rune) error {
	if len(s) == 0 || len(s) > len((*m)) {
		return ErrInvalidLength
	}

	z := len(s) - 1

	for y, char := range s {
		x, err := charToIndex(char)
		if err != nil {
			return err
		}

		(*m)[z][y][x] = true
	}

	return nil
}

func (m *CharMatrix) Contains(s []rune) bool {
	if len(s) == 0 || len(s) > len((*m)) {
		return false
	}

	z := len(s) - 1

	// Loop over all rows from the end to the start.
	for y := len(s) - 1; y >= 0; y-- {
		char := s[y]

		x, err := charToIndex(char)
		if err != nil || !(*m)[z][y][x] {
			return false
		}
	}

	return true
}

func (m *CharMatrix) Unset(s []rune) error {
	if len(s) == 0 || len(s) > len((*m)) {
		return ErrInvalidLength
	}

	z := len(s) - 1

	// Initially unset the last character in the input slice.
	x, err := charToIndex(s[z])
	if err != nil {
		return err
	}

	(*m)[z][z][x] = false

	// Loop over all rows from the end to the start, including the last one, excluding first.
	for y := len(s) - 1; y > 0; y-- {
		// Check if the current row has a set count of 0.
		if m.hasSetCount(z, y) > 0 {
			continue
		}

		// Convert the character in the row above to its index, if current row is empty.
		x, err := charToIndex(s[y-1])
		if err != nil {
			return err
		}

		// Unset the character in the row above if the current row has a set count of 0.
		(*m)[z][y-1][x] = false
	}

	return nil
}

func (m *CharMatrix) PrettyPrint(size int) {
	if size < 1 || size > len((*m)) {
		panic(ErrInvalidLength)
	}

	z := size - 1 // Corresponds to the depth of the matrix for the given size (zero-based indexing).
	maxIndexLength := m.calculateMaxIndexLength(z)

	{ // Print table header.
		headerPadding := fmt.Sprintf("%*s", maxIndexLength+2, " ")
		fmt.Print(headerPadding)

		for i := 0; i < totalCharactersCount; i++ {
			c, err := indexToChar(i)
			if err != nil {
				fmt.Printf("? ")

				continue
			}

			fmt.Printf("%c ", c)
		}

		fmt.Println()
	}

	// Iterate through each row in the 2D slice for the given size and print its contents.
	for y := 0; y < len((*m)[z]); y++ {
		// Align the row labels with the calculated padding.
		fmt.Printf("%*d: ", maxIndexLength, y)

		// Print each cell in the row, marking 'X' for character being at the position or '.'.
		for x := 0; x < len((*m)[z][y]); x++ {
			if (*m)[z][y][x] {
				fmt.Print("X ")
			} else {
				fmt.Print(". ")
			}
		}

		fmt.Println()
	}
}
